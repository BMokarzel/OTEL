package main

import (
	"context"
	"log"
	"os"

	"github.com/BMokarzel/OTEL/server-a/config"
	"github.com/BMokarzel/OTEL/server-a/internal/controller"
	"github.com/BMokarzel/OTEL/server-a/internal/gateway"
	"github.com/BMokarzel/OTEL/server-a/internal/service"
	"github.com/BMokarzel/OTEL/server-a/pkg/http"
	"github.com/BMokarzel/OTEL/server-a/pkg/logger"
	otl "github.com/BMokarzel/OTEL/server-a/pkg/otel"
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
		log.Fatalf("Error loading configs. Error %s", err)
		os.Exit(1)
	}
	/*
		err = otl.InitOtel(cfg.ServerName, cfg.CollectorUrl)
		if err != nil {
			log.Fatalf("Error to init otel. Error %s", err)
			os.Exit(1)
		}
	*/
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

	serverBHttpClient, err := http.NewHttpClient(cfg.ServerBUrl, logger, otel)
	if err != nil {
		log.Fatalf("Error to create new client http. Error %s", err)
		os.Exit(1)
	}

	serviceB := gateway.New(logger, serverBHttpClient, otel)

	service := service.New(logger, otel, serviceB)

	handler := controller.New(service, otel, server.Router)

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
		r.Post("/", a.handler.GetWeather)
		r.Get("/health", a.handler.HealthCheck)
	})
}
