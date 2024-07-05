package config_test

import (
	"context"
	"os"
	"testing"
	apperrors "weathersvc/app/app_errors"
	"weathersvc/app/config"

	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	os.Clearenv()
	ctx := context.Background()
	t.Run("Should not fail to create NewApp when config items are not missing", func(t *testing.T) {
		os.Setenv("WEATHER_ID", "fakeID")
		os.Setenv("WEATHER_HOST", "fakeHost")
		os.Setenv("PORT", "8081")
		os.Setenv("ENV", "testing")
		os.Setenv("SERVICE_URL", "fakevalue")
		expected := config.App{
			Port: "8081",
			Env:  "testing",
			WeatherClientConfig: config.WeatherClientConfig{
				Host:  "fakeHost",
				AppID: "fakeID",
			},
			ServiceURL: "fakevalue",
		}
		conf := config.NewAppConfig()
		resp, err := conf.NewApp(ctx)
		assert.NoError(t, err, "No errors expected for Config")
		assert.EqualValues(t, expected.WeatherClientConfig.AppID, resp.AppID)
		assert.EqualValues(t, expected.Port, resp.Port)
		assert.EqualValues(t, expected.Env, resp.Env)
		assert.EqualValues(t, resp.ServiceURL, expected.ServiceURL)
		assert.EqualValues(t, expected.WeatherClientConfig.AppID, resp.WeatherClientConfig.AppID)
		assert.EqualValues(t, expected.WeatherClientConfig.Host, resp.WeatherClientConfig.Host)
	})
	t.Run("Should fail to create NewApp when config Weather Host is missing", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("WEATHER_ID", "fakeID")
		os.Setenv("PORT", "8081")
		os.Setenv("ENV", "testing")
		os.Setenv("SERVICE_URL", "fakevalue")
		resp, err := config.NewAppConfig().NewApp(ctx)
		assert.EqualError(t, err, apperrors.CreateMissingConfigError(apperrors.ErrMissingConfig, "Weather Host").Error(), "Missing config for weather host")
		assert.Nil(t, resp)
	})
	t.Run("Should fail to create NewApp when config Weather AppID is missing", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("WEATHER_HOST", "fakeHost")
		os.Setenv("PORT", "8081")
		os.Setenv("ENV", "testing")
		os.Setenv("SERVICE_URL", "fakevalue")
		resp, err := config.NewAppConfig().NewApp(ctx)
		assert.EqualError(t, err, apperrors.CreateMissingConfigError(apperrors.ErrMissingConfig, "Weather App ID").Error(), "Missing config for weather host")
		assert.Nil(t, resp)
	})
	t.Run("Should not fail to create NewApp when config ServiceUrl is missing", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("WEATHER_ID", "fakeID")
		os.Setenv("WEATHER_HOST", "fakeHost")
		os.Setenv("PORT", "8081")
		os.Setenv("ENV", "testing")
		expected := config.App{
			Port: "8081",
			Env:  "testing",
			WeatherClientConfig: config.WeatherClientConfig{
				Host:  "fakeHost",
				AppID: "fakeID",
			},
			ServiceURL: "http://localhost",
		}
		resp, err := config.NewAppConfig().NewApp(ctx)
		assert.NoError(t, err, "No errors expected for Config")
		assert.EqualValues(t, expected.WeatherClientConfig.AppID, resp.AppID)
		assert.EqualValues(t, expected.WeatherClientConfig.Host, resp.Host)
		assert.EqualValues(t, expected.Port, resp.Port)
		assert.EqualValues(t, expected.Env, resp.Env)
		assert.EqualValues(t, resp.ServiceURL, expected.ServiceURL)
		assert.EqualValues(t, expected.WeatherClientConfig.AppID, resp.WeatherClientConfig.AppID)
		assert.EqualValues(t, expected.WeatherClientConfig.Host, resp.WeatherClientConfig.Host)
	})
	t.Run("Should not fail to create NewApp when config Port is missing", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("WEATHER_ID", "fakeID")
		os.Setenv("WEATHER_HOST", "fakeHost")
		os.Setenv("ENV", "testing")
		os.Setenv("SERVICE_URL", "fakevalue")
		expected := config.App{
			Port: "8080",
			Env:  "testing",
			WeatherClientConfig: config.WeatherClientConfig{
				Host:  "fakeHost",
				AppID: "fakeID",
			},
			ServiceURL: "fakevalue",
		}
		resp, err := config.NewAppConfig().NewApp(ctx)
		assert.NoError(t, err, "No errors expected for Config")
		assert.EqualValues(t, expected.WeatherClientConfig.AppID, resp.AppID)
		assert.EqualValues(t, expected.Port, resp.Port)
		assert.EqualValues(t, expected.Env, resp.Env)
		assert.EqualValues(t, resp.ServiceURL, expected.ServiceURL)
		assert.EqualValues(t, expected.WeatherClientConfig.AppID, resp.WeatherClientConfig.AppID)
		assert.EqualValues(t, expected.WeatherClientConfig.Host, resp.WeatherClientConfig.Host)
	})
}
