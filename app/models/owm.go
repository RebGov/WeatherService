/*
owm.go: This model is based on response as provided by open weather map api.
*/
package models

// WeatherResponse represents the structure of the JSON response
type WeatherResponse struct {
	Coord      Coord              `json:"coord"`
	Weather    []Weather          `json:"weather"`
	Base       string             `json:"base"`
	Main       Main               `json:"main"`
	Visibility int                `json:"visibility"`
	Wind       Wind               `json:"wind"`
	Rain       map[string]float64 `json:"rain"`
	Clouds     Clouds             `json:"clouds"`
	Dt         int64              `json:"dt"`
	Sys        Sys                `json:"sys"`
	Timezone   int                `json:"timezone"`
	ID         int                `json:"id"`
	Name       string             `json:"name"`
	Cod        int                `json:"cod"`
}

// Coord represents the coordinates
type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// Weather represents weather information
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Main represents main weather information
type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}

// Wind represents wind information
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

// Clouds represents cloud information
type Clouds struct {
	All int `json:"all"`
}

// Sys represents system information
type Sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int64  `json:"sunrise"`
	Sunset  int64  `json:"sunset"`
}
