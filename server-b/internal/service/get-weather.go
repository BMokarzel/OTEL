package service

import (
	"context"

	controller_dto "github.com/BMokarzel/OTEL/server-b/internal/controller/dto"
	viacep "github.com/BMokarzel/OTEL/server-b/internal/gateway/via-cep"
	"github.com/BMokarzel/OTEL/server-b/pkg/errors"
)

func (s *Service) GetWeather(ctx context.Context, zipCode string) (controller_dto.GetWeatherOutput, error) {
	ctx, mainSpan := s.Otel.OTELTracer.Start(ctx, "service.GetWather")
	defer mainSpan.End()

	var viaCepRes viacep.ViaCepOutput

	viaCepRes, err := s.ViaCep.GetLocation(ctx, zipCode)
	if err != nil {
		return controller_dto.GetWeatherOutput{}, errors.NewNotFoundError("can not find zipcode")
	}

	watherRes, err := s.WeatherApi.GetWeather(ctx, viaCepRes.Location)
	if err != nil {
		return controller_dto.GetWeatherOutput{}, errors.NewNotFoundError("error to get real time weather")
	}

	{
		_, span := s.Otel.OTELTracer.Start(ctx, "converting-temperatures")
		defer span.End()

		tempF := watherRes.Weather.TempC*1.8 + 32

		tempK := watherRes.Weather.TempC + 273

		response := controller_dto.GetWeatherOutput{
			City:  viaCepRes.Location,
			TempC: watherRes.Weather.TempC,
			TempF: tempF,
			TempK: tempK,
		}
		return response, nil
	}
}
