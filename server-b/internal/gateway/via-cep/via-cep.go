package viacep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BMokarzel/OTEL/server-b/pkg/logger"
	"github.com/BMokarzel/OTEL/server-b/pkg/otel"
)

type ViaCep struct {
	logger *logger.Logger
	URL    string
	Otel   *otel.Otel
}

type ViaCepOutput struct {
	Location string `json:"localidade"`
	Error    string `json:"erro"`
}

func New(logger *logger.Logger, url string, otel *otel.Otel) *ViaCep {
	return &ViaCep{
		logger: logger,
		URL:    url,
		Otel:   otel,
	}
}

func (v *ViaCep) GetLocation(ctx context.Context, cep string) (ViaCepOutput, error) {

	url := fmt.Sprintf("%s/%s/json/", v.URL, cep)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ViaCepOutput{}, err
	}

	v.logger.Info(ctx).Msg("[DEBUG] Request: ", req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ViaCepOutput{}, err
	}

	if res.StatusCode > 299 {

		v.logger.Info(ctx).Msg("[DEBUG] Response: ", res)

		return ViaCepOutput{}, fmt.Errorf("")

	} else {

		var response ViaCepOutput

		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return ViaCepOutput{}, err
		}

		v.logger.Info(ctx).Msg("[DEBUG] Response: ", res)

		return response, nil
	}
}
