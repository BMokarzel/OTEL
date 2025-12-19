package weatherapi

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/BMokarzel/OTEL/server-b/pkg/errors"
	"github.com/BMokarzel/OTEL/server-b/pkg/logger"
	"github.com/BMokarzel/OTEL/server-b/pkg/otel"
)

type WeatherApi struct {
	logger *logger.Logger
	URL    string
	Key    string
	Otel   *otel.Otel
}

type GetWeatherOutput struct {
	Location Location
	Weather  Current `json:"current"`
}

type Current struct {
	TempC float64 `json:"temp_c"`
}

type Location struct {
	Name      string `json:"name"`
	Region    string `json:"region"`
	Country   string `json:"country"`
	LocalTime string `json:"localtime"`
}

func New(logger *logger.Logger, url, key string, otel *otel.Otel) *WeatherApi {
	return &WeatherApi{
		logger: logger,
		URL:    url,
		Key:    key,
		Otel:   otel,
	}
}

func (k *WeatherApi) GetWeather(ctx context.Context, location string) (GetWeatherOutput, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, k.URL, nil)
	if err != nil {
		return GetWeatherOutput{}, err
	}

	c := req.URL.Query()
	c.Add("key", k.Key)
	c.Add("q", location)
	req.URL.RawQuery = c.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return GetWeatherOutput{}, err
	}

	defer res.Body.Close()

	var response GetWeatherOutput

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return GetWeatherOutput{}, err
	}

	switch {
	case res.StatusCode < 300:
		return response, nil
	case res.StatusCode == 400:
		return GetWeatherOutput{}, errors.NewBadRequestError("")
	case res.StatusCode == 404:
		return GetWeatherOutput{}, errors.NewNotFoundError("")
	case res.StatusCode == 422:
		return GetWeatherOutput{}, errors.NewUnprocessableError("")
	default:
		k.logger.Error(ctx).Msg("Unkown status code. Code: %d", res.StatusCode)
		return GetWeatherOutput{}, errors.NewInternalServerError("")
	}
}
