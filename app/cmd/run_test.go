package cmd_test

import (
	"context"
	"errors"
	"os"
	"testing"

	cmd "weathersvc/app/cmd"

	"github.com/stretchr/testify/assert"
)

func TestCMD_Execute(t *testing.T) {
	ctx := context.Background()
	os.Setenv("WEATHER_ID", "fakeID")
	os.Setenv("WEATHER_HOST", "https://api.openweathermap.org/data/2.5/weather")
	os.Setenv("PORT", "8081")
	os.Setenv("ENV", "testing")
	os.Setenv("SERVICE_URL", "http://fakevalue.com/")
	t.Run("Should fail when config is invalid", func(t *testing.T) {
		gotErr := cmd.Execute(ctx)
		assert.EqualError(t, gotErr, errors.New("config `WEATHER_ID` is invalid").Error())
	})
	t.Run("Should fail to run due to missing `Weather App ID` config", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("WEATHER_HOST", "fakeHost")
		os.Setenv("PORT", "8081")
		os.Setenv("ENV", "testing")
		os.Setenv("SERVICE_URL", "fakevalue")
		gotErr := cmd.Execute(ctx)
		assert.EqualError(t, gotErr, errors.New("failed to start service: missing required config for `Weather App ID`").Error())
	})
}
