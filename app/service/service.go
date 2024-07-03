package service

import (
	"context"
	"weathersvc/app/config"
	openweather "weathersvc/app/open_weather"
)

type WeatherService struct {
	config        *config.App
	WeatherClient *openweather.Client
}

func NewService(ctx context.Context, conf *config.App) (*WeatherService, error) {
	client := openweather.NewClient(conf)
	err := client.ApiTest()
	if err != nil {
		return nil, err
	}
	return &WeatherService{
		config:        conf,
		WeatherClient: client,
	}, nil
}
