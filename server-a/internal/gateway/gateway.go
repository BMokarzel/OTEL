package gateway

import (
	"github.com/BMokarzel/OTEL/server-a/pkg/http"
	"github.com/BMokarzel/OTEL/server-a/pkg/logger"
	"github.com/BMokarzel/OTEL/server-a/pkg/otel"
)

type ServiceB struct {
	logger *logger.Logger
	Client *http.HttpClient
	Otel   *otel.Otel
}

func New(logger *logger.Logger, client *http.HttpClient, otel *otel.Otel) *ServiceB {
	return &ServiceB{
		logger: logger,
		Client: client,
		Otel:   otel,
	}
}
