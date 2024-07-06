package service

import (
	"context"
	"weathersvc/app/config"
	openweather "weathersvc/app/open_weather"
)

type Service interface {
	// GetWeather ctx, latitude, longitude
	GetWeather(ctx context.Context, lat, lon float64) (WeatherCond, error)
	ValidateSvc(ctx context.Context) error
}
type service struct {
	Config        *config.App
	WeatherClient openweather.Client
}

func NewService(ctx context.Context, conf *config.App) Service {
	cl := openweather.NewClient(conf)
	return &service{
		Config:        conf,
		WeatherClient: cl,
	}
}

func (s *service) ValidateSvc(ctx context.Context) error {
	return s.WeatherClient.ApiTest()
}
