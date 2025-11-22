package service

import (
	"context"

	controller_dto "github.com/BMokarzel/OTEL/server-b/internal/controller/dto"
	viacep "github.com/BMokarzel/OTEL/server-b/internal/gateway/via-cep"
)

func (s *Service) GetWeather(ctx context.Context, zipCode string) (interface{}, int) {

	ctx, mainSpan := s.Otel.OTELTracer.Start(ctx, "service.GetWather")
	defer mainSpan.End()

	var viaCepRes viacep.ViaCepOutput

	viaCepRes, err := s.ViaCep.GetLocation(ctx, zipCode)
	if err != nil || viaCepRes.Error == "true" {
		return controller_dto.ErrorOutput{
			Message: "can not find zipcode",
		}, 404
	}

	watherRes, err := s.WeatherApi.GetWeather(ctx, viaCepRes.Location)
	if err != nil {
		return controller_dto.ErrorOutput{
			Message: "problem to get real time weather. If the problem persists, contact support",
		}, 422
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
		return response, 200
	}
}
