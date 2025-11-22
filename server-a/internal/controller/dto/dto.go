package controller_dto

type GetWeatherInput struct {
	ZipCode string `json:"cep"`
}

type GetWeatherOutput struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ErrorOutput struct {
	Message string `json:"error"`
}
