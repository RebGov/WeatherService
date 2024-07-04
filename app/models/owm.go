/*
owm.go: This model is based on response as provided by open weather map api.
Only the items uttilizled in this service are maintained
*/
package models

// WeatherResponse represents the structure of the JSON response
type WeatherResponse struct {
	Weather []Weather `json:"weather"`
	Main    Main      `json:"main"`
	Wind    Wind      `json:"wind"`
	Cod     int       `json:"cod"`
}

type Weather struct {
	Description string `json:"description"`
}

type Main struct {
	FeelsLike float64 `json:"feels_like"`
}

type Wind struct {
	Speed float64 `json:"speed"`
}
