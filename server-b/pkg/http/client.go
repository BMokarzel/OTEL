package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/BMokarzel/OTEL/server-b/pkg/errors"
	"github.com/BMokarzel/OTEL/server-b/pkg/logger"
	otel_pkg "github.com/BMokarzel/OTEL/server-b/pkg/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type HttpClientOps struct {
	URL string
}

type HttpClient struct {
	URL        string
	logger     *logger.Logger
	HttpClient http.Client
	Otel       *otel_pkg.Otel
}

type Params struct {
	Path   string
	Method string
	Header map[string]string
	Body   *bytes.Reader
}

func NewHttpClient(url string, logger *logger.Logger, otel *otel_pkg.Otel) (*HttpClient, error) {

	client := http.Client{
		Timeout: time.Minute * 3,
	}

	return &HttpClient{
		URL:        url,
		logger:     logger,
		HttpClient: client,
		Otel:       otel,
	}, nil
}

func (c *HttpClient) Call(ctx context.Context, param Params, response interface{}) error {

	ctx, span := c.Otel.OTELTracer.Start(ctx, "http.Call")

	defer span.End()

	var req *http.Request
	var err error

	if param.Body != nil {
		req, err = http.NewRequestWithContext(ctx, param.Method, param.Path, param.Body)
		if err != nil {
			span.RecordError(err)
			return err
		}
	} else {
		req, err = http.NewRequestWithContext(ctx, param.Method, param.Path, nil)
		if err != nil {
			span.RecordError(err)
			return err
		}
	}

	dumpReq, err := httputil.DumpRequest(req, true)
	if err == nil {
		c.logger.Info(ctx).Msg("[DEBUG] Client request - %s", dumpReq)
	}

	carrier := propagation.HeaderCarrier(req.Header)

	otel.GetTextMapPropagator().Inject(ctx, carrier)

	for k, v := range param.Header {
		req.Header.Set(k, v)
	}

	//adicionar no header os termos para o otel

	res, err := c.HttpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		return err
	}

	if res.StatusCode > 299 {

		var genericError errors.GenericError

		err := json.NewDecoder(res.Body).Decode(&genericError)
		if err != nil {
			span.RecordError(err)
			return err
		}

		switch res.StatusCode {
		case 404:
			e := errors.NewBadRequestError(genericError.Message)
			span.RecordError(e)
			return e
		case 422:
			e := errors.NewUnprocessableError(genericError.Message)
			span.RecordError(e)
			return e
		default:
			e := errors.NewBadRequestError(genericError.Message)
			span.RecordError(e)
			return e
		}

	} else {

		err = json.NewDecoder(res.Body).Decode(response)
		if err != nil {
			span.RecordError(err)
			return err
		}

	}

	return nil
}
