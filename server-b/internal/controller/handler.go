package controller

import (
	"encoding/json"
	"net/http"

	"github.com/BMokarzel/OTEL/server-b/internal/service"
	"github.com/BMokarzel/OTEL/server-b/pkg/errors"
	pkg_http "github.com/BMokarzel/OTEL/server-b/pkg/http"
	"github.com/BMokarzel/OTEL/server-b/pkg/logger"
	otl "github.com/BMokarzel/OTEL/server-b/pkg/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Handler struct {
	logger  *logger.Logger
	Service *service.Service
	Otel    *otl.Otel
	Router  *pkg_http.Router
}

func New(logger *logger.Logger, service *service.Service, otel *otl.Otel) *Handler {
	return &Handler{
		logger:  logger,
		Service: service,
		Otel:    otel,
	}
}

func (h *Handler) GetWeather(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := h.Otel.OTELTracer.Start(ctx, "[GET] /")
	defer span.End()

	zipCode := r.URL.Query().Get("zipCode")

	res, err := h.Service.GetWeather(ctx, zipCode)
	if err != nil {
		h.Router.ErrorHandler(w, r, err)
		return
	}

	body, err := json.Marshal(res)
	if err != nil {
		h.logger.Error(ctx).Msg("Error to parse response body. Error: %s", err)
		span.RecordError(err)
		h.Router.ErrorHandler(w, r, errors.NewInternalServerError(""))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
