package service

import (
	"context"
	"fmt"
)

type WeatherCond struct {
	Temp      Temperature
	Condition string
	Wind      string
}
type Temperature string

var (
	UnknownTemp     Temperature = "unknown"
	ExtremeFreezing Temperature = "extreme freezing"
	SubFreezing     Temperature = "sub-freezing"
	Freezing        Temperature = "freezing"
	Cold            Temperature = "cold"
	Moderate        Temperature = "moderate"
	Hot             Temperature = "hot"
	ExtremeHot      Temperature = "extremely hot"
)

func (w *WeatherService) GetWeather(ctx context.Context, lat, lon float64) (WeatherCond, error) {
	resp, err := w.WeatherClient.GetWeather(fmt.Sprintf("%f", lat), fmt.Sprintf("%f", lon))
	if err != nil {
		return WeatherCond{}, err
	}
	tempCond := w.buildTempCondition(resp.Main.FeelsLike)

	if len(resp.Weather) == 0 {
		return WeatherCond{
			Temp:      tempCond,
			Condition: "unknown",
			Wind:      fmt.Sprintf("%.2f miles/hour", resp.Wind.Speed),
		}, nil
	}
	return WeatherCond{
		Temp:      tempCond,
		Condition: resp.Weather[0].Description,
		Wind:      fmt.Sprintf("%.2f miles/hour", resp.Wind.Speed),
	}, nil
}

func (w *WeatherService) buildTempCondition(temp float64) Temperature {
	switch {
	case temp <= -10:
		return ExtremeFreezing
	case temp > -10 && temp <= 20:
		return SubFreezing
	case temp > 21 && temp <= 32:
		return Freezing
	case temp > 32 && temp <= 65:
		return Cold
	case temp > 65 && temp <= 75:
		return Moderate
	case temp > 75 && temp <= 89:
		return Hot
	case temp > 90:
		return ExtremeHot
	default:
		return UnknownTemp
	}
}
