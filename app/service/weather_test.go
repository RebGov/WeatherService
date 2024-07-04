package service

import (
	"testing"
	"weathersvc/app/config"
	openweather "weathersvc/app/open_weather"

	"github.com/stretchr/testify/assert"
)

func TestWeatherService_buildTempCondition(t *testing.T) {
	svc := &WeatherService{
		config:        &config.App{},
		WeatherClient: &openweather.Client{},
	}
	t.Run("Should return `extremely hot`", func(t *testing.T) {
		got := svc.buildTempCondition(95.0)
		assert.EqualValues(t, got, ExtremeHot)
	})
	t.Run("Should return `hot`", func(t *testing.T) {
		got := svc.buildTempCondition(76.0)
		assert.EqualValues(t, got, Hot)
	})
	t.Run("Should return `moderate`", func(t *testing.T) {
		got := svc.buildTempCondition(71.0)
		assert.EqualValues(t, got, Moderate)
	})
	t.Run("Should return `Cold`", func(t *testing.T) {
		got := svc.buildTempCondition(47.0)
		assert.EqualValues(t, got, Cold)
	})
	t.Run("Should return `Freezing`", func(t *testing.T) {
		got := svc.buildTempCondition(32.0)
		assert.EqualValues(t, got, Freezing)
	})
	t.Run("Should return `Sub-Freezing`", func(t *testing.T) {
		got := svc.buildTempCondition(-9.0)
		assert.EqualValues(t, got, SubFreezing)
	})
	t.Run("Should return `Extreme Freezing`", func(t *testing.T) {
		got := svc.buildTempCondition(-15.0)
		assert.EqualValues(t, got, ExtremeFreezing)
	})
}
