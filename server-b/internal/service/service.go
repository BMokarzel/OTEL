package service

import (
	viacep "github.com/BMokarzel/OTEL/server-b/internal/gateway/via-cep"
	weatherapi "github.com/BMokarzel/OTEL/server-b/internal/gateway/weather-api"
	"github.com/BMokarzel/OTEL/server-b/pkg/logger"
	"github.com/BMokarzel/OTEL/server-b/pkg/otel"
)

type Service struct {
	logger     *logger.Logger
	ViaCep     *viacep.ViaCep
	WeatherApi *weatherapi.WeatherApi
	Otel       *otel.Otel
}

func New(logger *logger.Logger, viaCep *viacep.ViaCep, weatherApi *weatherapi.WeatherApi, otel *otel.Otel) *Service {
	return &Service{
		logger:     logger,
		ViaCep:     viaCep,
		WeatherApi: weatherApi,
		Otel:       otel,
	}
}
