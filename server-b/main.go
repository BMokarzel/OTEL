package main

import (
	"context"
	"log"
	"os"

	"github.com/BMokarzel/OTEL/server-b/config"
	"github.com/BMokarzel/OTEL/server-b/internal/controller"
	viacep "github.com/BMokarzel/OTEL/server-b/internal/gateway/via-cep"
	weatherapi "github.com/BMokarzel/OTEL/server-b/internal/gateway/weather-api"
	"github.com/BMokarzel/OTEL/server-b/internal/service"
	"github.com/BMokarzel/OTEL/server-b/pkg/http"
	"github.com/BMokarzel/OTEL/server-b/pkg/logger"
	otl "github.com/BMokarzel/OTEL/server-b/pkg/otel"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
)

type api struct {
	handler *controller.Handler
	router  *http.Router
}

func main() {

	cfg, err := config.LoadConfigs()
	if err != nil {
		log.Fatal("error to load configs. Error: %s", err)
		os.Exit(1)
	}

	logger := logger.New(cfg.Build, cfg.ServerName)

	serverOpts := http.ServerOpts{
		ServerPort:          cfg.ServerPort,
		ServerVersion:       cfg.Build,
		ServerEnableTracing: cfg.ServerEnableTracing,
		ServerEnv:           cfg.ServerEnv,
		ServerName:          cfg.ServerName,
	}

	otelOpts := http.ServerOtel{
		ServerName:   cfg.ServerName,
		CollectorUrl: cfg.CollectorUrl,
	}

	server := http.NewServer(logger, serverOpts, otelOpts)

	otel := otl.New(otel.Tracer(cfg.ServerName), cfg.ServerName)

	viaCep := viacep.New(logger, cfg.ViaCepUrl, otel)

	weatherapi := weatherapi.New(logger, cfg.WeatherApiUrl, cfg.WeatherApiKey, otel)

	service := service.New(logger, viaCep, weatherapi, otel)

	handler := controller.New(logger, service, otel)

	api := api{
		handler: handler,
		router:  server.Router,
	}

	api.Routes()

	if err := server.Init(); err != nil {
		logger.Error(context.Background()).Msg("error to run server: %v", err)
		os.Exit(1)
	}
}

func (a *api) Routes() {
	a.router.Engine.Route("/", func(r chi.Router) {
		r.Get("/", a.handler.GetWeather)
	})
}
