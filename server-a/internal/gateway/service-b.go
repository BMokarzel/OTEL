package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	http_package "github.com/BMokarzel/OTEL/server-a/pkg/http"
)

func (s *ServiceB) Call(ctx context.Context, zipCode string) (json.RawMessage, error) {

	ctx, span := s.Otel.OTELTracer.Start(ctx, "server-b.Call")

	defer span.End()

	params := http_package.Params{
		Path:   fmt.Sprintf("%s", s.Client.URL),
		Method: http.MethodGet,
		Header: map[string]string{
			"Content-Type": "application/json",
		},
		Query: map[string]string{
			"zipCode": zipCode,
		},
		Body: nil,
	}

	response := new(json.RawMessage)

	err := s.Client.Call(ctx, params, response)
	if err != nil {
		s.logger.Error(ctx).Msg("%s", err)
		return nil, err
	}

	return *response, nil

}
