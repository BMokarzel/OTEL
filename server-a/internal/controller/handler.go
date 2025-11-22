package controller

import (
	"encoding/json"
	"net/http"
	"time"

	controller_dto "github.com/BMokarzel/OTEL/server-a/internal/controller/dto"
	"github.com/BMokarzel/OTEL/server-a/internal/service"
	http_package "github.com/BMokarzel/OTEL/server-a/pkg/http"
	otl "github.com/BMokarzel/OTEL/server-a/pkg/otel"
	"github.com/BMokarzel/OTEL/server-b/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Handler struct {
	logger  *logger.Logger
	Service *service.Service
	Otel    *otl.Otel
	Router  *http_package.Router
}

func New(service *service.Service, otel *otl.Otel, router *http_package.Router) *Handler {
	return &Handler{
		Service: service,
		Otel:    otel,
		Router:  router,
	}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := h.Otel.OTELTracer.Start(ctx, "[GET] /")
	defer span.End()

	time.Sleep(time.Second * 2)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetWeather(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := h.Otel.OTELTracer.Start(ctx, "[POST] /")
	defer span.End()

	var input controller_dto.GetWeatherInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request body"))
		return
	}

	res, err := h.Service.GetWeather(ctx, input)
	if err != nil {
		h.Router.ErrorHandler(w, r, err)
		return
	}

	body, err := json.Marshal(res)
	if err != nil {
		h.logger.Error(ctx).Msg("error to parse response body. Error: %s", err)
		span.RecordError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
