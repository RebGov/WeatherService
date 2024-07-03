package service

import (
	"context"
	"fmt"
)

func convertFloatToString(f float32) string {
	return fmt.Sprintf("%f", f)
}

func (w *WeatherService) GetWeather(ctx context.Context, lat, long float32) error {
	resp, err := w.WeatherClient.GetWeather(convertFloatToString(lat), convertFloatToString(long))
	if err != nil {
		return err
	}
	w.buildCondition(resp.Main.FeelsLike)
	return nil
}

func (w *WeatherService) buildCondition(temp float64) {
	switch {
	// case temp < :

	}
}
