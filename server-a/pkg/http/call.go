package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httputil"

	"github.com/BMokarzel/OTEL/server-a/pkg/errors"
	"github.com/BMokarzel/OTEL/server-a/pkg/logger"
	otl "github.com/BMokarzel/OTEL/server-a/pkg/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Param struct {
	Path   string
	Method string
	Header map[string]string
	Body   *bytes.Reader
}

type Cl struct {
	Client *http.Client
	logger *logger.Logger
	Otel   *otl.Otel
}

func (c *Cl) Call(ctx context.Context, param Param, response interface{}) error {
	ctx, span := c.Otel.OTELTracer.Start(ctx, "http.Call")
	defer span.End()

	var req *http.Request
	var err error

	if param.Body != nil {
		req, err = http.NewRequestWithContext(ctx, param.Method, param.Path, param.Body)
		if err != nil {
			c.logger.Error(ctx).Msg("error to create new request with context and body. Error: %s", err)
			span.RecordError(err)
			return err
		}
	} else {
		req, err = http.NewRequestWithContext(ctx, param.Method, param.Path, nil)
		if err != nil {
			c.logger.Error(ctx).Msg("error to create new request with context. Error: %s", err)
			span.RecordError(err)
			return err
		}
	}

	dumpReq, err := httputil.DumpRequest(req, true)
	if err == nil {
		c.logger.Info(ctx).Msg("[DEBUG] Client request - %s", dumpReq)
	}

	for k, v := range param.Header {
		req.Header.Set(k, v)
	}

	otelPropagartor := otel.GetTextMapPropagator()

	otelPropagartor.Inject(ctx, propagation.HeaderCarrier(req.Header))

	res, err := c.Client.Do(req)
	if err != nil {
		c.logger.Error(ctx).Msg("error to send request. Error: %s", err)
		span.RecordError(err)
		return err
	}

	if res.StatusCode > 299 {

		var genericError errors.GenericError

		err := json.NewDecoder(res.Body).Decode(&genericError)
		if err != nil {
			c.logger.Error(ctx).Msg("error to parse response body. Error: %s", err)
			span.RecordError(err)
			return err
		}

		switch res.StatusCode {
		case 404:
			return errors.NewBadRequestError(genericError.Message)
		case 422:
			return errors.NewUnprocessableError(genericError.Message)
		default:
			return errors.NewBadRequestError(genericError.Message)
		}

	} else {

		err = json.NewDecoder(res.Body).Decode(response)
		if err != nil {
			c.logger.Error(ctx).Msg("error to parse responde body. Error: %s", err)
			span.RecordError(err)
			return err
		}

		return nil
	}
}
