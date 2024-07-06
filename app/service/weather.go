package service

import (
	"context"
	"fmt"
	"log"
)

type WeatherCond struct {
	Temp      Temperature
	Condition string
	Wind      Wind
}
type Temperature string
type Wind string

const (
	UnknownTemp    Temperature = "unknown"
	subFreezing    Temperature = "sub-freezing"
	Freezing       Temperature = "freezing"
	cold           Temperature = "cold"
	moderate       Temperature = "moderate"
	warm           Temperature = "warm"
	hot            Temperature = "hot"
	extremeHot     Temperature = "extremely hot"
	unknownWind    Wind        = "unknown wind"
	calm           Wind        = "calm winds"
	lightAir       Wind        = "light air"
	lightBreeze    Wind        = "light breeze"
	gentalBreeze   Wind        = "gentle breeze"
	moderateBreeze Wind        = "moderate breeze"
	freshBreeze    Wind        = "fresh breeze"
	strongBreeze   Wind        = "strong breeze"
	nearGale       Wind        = "near gale winds"
	gale           Wind        = "gale winds"
	severeGale     Wind        = "severe gale winds"
	storm          Wind        = "storm winds"
	violentStorm   Wind        = "violent storm winds"
	hurricane      Wind        = "hurricane/tornado winds"
)

// GetWeather ctx, latitude, longitude
func (w *service) GetWeather(ctx context.Context, lat, lon float64) (WeatherCond, error) {
	log.Println("AM HERE")
	sLat := fmt.Sprintf("%f", lat)
	sLon := fmt.Sprintf("%f", lon)
	log.Printf("Lat: %s, Lon:%s", sLat, sLon)
	resp, err := w.WeatherClient.GetWeather(sLat, sLon)
	if err != nil {
		return WeatherCond{}, err
	}
	tempCond := w.buildTempCondition(resp.Main.FeelsLike)
	windCond := w.buildWindCondition(resp.Wind.Speed)
	if len(resp.Weather) == 0 {
		return WeatherCond{
			Temp:      tempCond,
			Condition: "unknown",
			Wind:      windCond,
		}, nil
	}
	return WeatherCond{
		Temp:      tempCond,
		Condition: resp.Weather[0].Description,
		Wind:      windCond,
	}, nil
}

func (w *service) buildTempCondition(temp float64) Temperature {
	switch {
	case temp <= 32:
		return subFreezing
	case temp > 32 && temp <= 40:
		return Freezing
	case temp > 40 && temp <= 60:
		return cold
	case temp > 60 && temp <= 75:
		return moderate
	case temp > 75 && temp <= 90:
		return warm
	case temp > 90 && temp < 100:
		return hot
	case temp >= 100:
		return extremeHot
	default:
		return UnknownTemp
	}
}

func (w *service) buildWindCondition(s float64) Wind {
	switch {
	case s <= 0:
		return calm
	case s > 0 && s < 4:
		return lightAir
	case s >= 4 && s < 8:
		return lightBreeze
	case s >= 8 && s < 13:
		return gentalBreeze
	case s >= 13 && s < 19:
		return moderateBreeze
	case s >= 19 && s < 25:
		return freshBreeze
	case s >= 25 && s < 32:
		return strongBreeze
	case s >= 32 && s < 39:
		return nearGale
	case s >= 39 && s < 47:
		return gale
	case s >= 47 && s < 55:
		return severeGale
	case s >= 55 && s < 64:
		return storm
	case s >= 64 && s < 73:
		return violentStorm
	case s >= 73:
		return hurricane
	default:
		return unknownWind
	}
}
