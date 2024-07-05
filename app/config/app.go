package config

import (
	"context"
	"os"
	appErr "weathersvc/app/app_errors"
)

type AppConfig interface {
	NewApp(ctx context.Context) (*App, error)
}

type App struct {
	Port string
	Env  string
	WeatherClientConfig
	ServiceURL string
}

type WeatherClientConfig struct {
	Host  string
	AppID string
}

type appConfigImpl struct{}

func NewAppConfig() AppConfig {
	return &appConfigImpl{}
}

func (a *appConfigImpl) NewApp(ctx context.Context) (*App, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	wAppID := os.Getenv("WEATHER_ID")
	if wAppID == "" {
		return nil, appErr.CreateMissingConfigError(appErr.ErrMissingConfig, "Weather App ID")
	}
	wHost := os.Getenv("WEATHER_HOST")
	if wHost == "" {
		return nil, appErr.CreateMissingConfigError(appErr.ErrMissingConfig, "Weather Host")
	}
	svcUrl := os.Getenv("SERVICE_URL")
	if svcUrl == "" {
		svcUrl = "http://localhost"
	}
	return &App{
		Port: port,
		Env:  os.Getenv("ENV"),
		WeatherClientConfig: WeatherClientConfig{
			Host:  wHost,
			AppID: wAppID,
		},
		ServiceURL: svcUrl,
	}, nil
}
