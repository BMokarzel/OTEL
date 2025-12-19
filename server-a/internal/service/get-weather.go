package service

import (
	"context"
	"encoding/json"
	"regexp"

	controller_dto "github.com/BMokarzel/OTEL/server-a/internal/controller/dto"
	"github.com/BMokarzel/OTEL/server-a/pkg/errors"
)

func (s *Service) GetWeather(ctx context.Context, input controller_dto.GetWeatherInput) (json.RawMessage, error) {
	ctx, mainSpan := s.Otel.OTELTracer.Start(ctx, "service.GetWeather")

	defer mainSpan.End()

	{
		_, span := s.Otel.OTELTracer.Start(ctx, "validate.ZipCode")

		defer span.End()

		regex := regexp.MustCompile(`^\d{8}$`)

		if !regex.MatchString(input.ZipCode) {
			e := errors.NewUnprocessableError("invalid zipcode")
			s.logger.Error(ctx).Msg("invalid zipcode type")
			span.RecordError(e)
			return nil, e
		}
	}

	{
		res, err := s.ServiceB.Call(ctx, input.ZipCode)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}
