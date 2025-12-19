package viacep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BMokarzel/OTEL/server-b/pkg/errors"
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

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		v.logger.Error(ctx).Msg("Error to do request. Error: %s", err)
		return ViaCepOutput{}, err
	}

	defer res.Body.Close()

	var response ViaCepOutput

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		v.logger.Error(ctx).Msg("Error to parse request body. Error: %s", err)
		return ViaCepOutput{}, err
	}

	v.logger.Info(ctx).Msg("[DEBUG] Status code: %d. Message: %v", res.StatusCode, res.Body)

	switch {
	case response.Error != "":
		return ViaCepOutput{}, errors.NewBadRequestError("")
	case res.StatusCode < 300:
		return response, nil
	case res.StatusCode == 400:
		return ViaCepOutput{}, errors.NewBadRequestError("")
	case res.StatusCode == 404:
		return ViaCepOutput{}, errors.NewNotFoundError("")
	case res.StatusCode == 422:
		return ViaCepOutput{}, errors.NewUnprocessableError("")
	default:
		return ViaCepOutput{}, errors.NewInternalServerError("")
	}
}
