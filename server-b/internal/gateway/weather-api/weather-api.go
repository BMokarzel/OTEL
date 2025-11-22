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

	k.logger.Info(ctx).Msg("[DEBUG] Request: ", req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return GetWeatherOutput{}, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		k.logger.Info(ctx).Msg("[DEBUG] Response: ", res)

		return GetWeatherOutput{}, errors.NewNotFoundError("")

	} else if res.StatusCode > 299 && res.StatusCode != 404 {
		k.logger.Info(ctx).Msg("[DEBUG] Response: ", res)

		return GetWeatherOutput{}, errors.NewUnprocessableError("")
	} else {
		k.logger.Info(ctx).Msg("[DEBUG] Response: ", res)

		var response GetWeatherOutput

		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return GetWeatherOutput{}, err

		}

		return response, nil
	}

}
