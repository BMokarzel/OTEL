package service

import (
	"github.com/BMokarzel/OTEL/server-a/internal/gateway"
	"github.com/BMokarzel/OTEL/server-a/pkg/logger"
	"github.com/BMokarzel/OTEL/server-a/pkg/otel"
)

type Service struct {
	logger   *logger.Logger
	ServiceB *gateway.ServiceB
	Otel     *otel.Otel
}

func New(logger *logger.Logger, otel *otel.Otel, serviceB *gateway.ServiceB) *Service {
	return &Service{
		logger:   logger,
		ServiceB: serviceB,
		Otel:     otel,
	}
}
