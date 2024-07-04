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
		got := svc.buildTempCondition(101.0)
		assert.EqualValues(t, got, extremeHot)
	})
	t.Run("Should return `hot`", func(t *testing.T) {
		got := svc.buildTempCondition(92.0)
		assert.EqualValues(t, got, hot)
	})
	t.Run("Should return `moderate`", func(t *testing.T) {
		got := svc.buildTempCondition(69.0)
		assert.EqualValues(t, got, moderate)
	})
	t.Run("Should return `Cold`", func(t *testing.T) {
		got := svc.buildTempCondition(47.0)
		assert.EqualValues(t, got, cold)
	})
	t.Run("Should return `Freezing`", func(t *testing.T) {
		got := svc.buildTempCondition(33.0)
		assert.EqualValues(t, got, Freezing)
	})
	t.Run("Should return `Sub-Freezing`", func(t *testing.T) {
		got := svc.buildTempCondition(-9.0)
		assert.EqualValues(t, got, subFreezing)
	})
	t.Run("Should return `Warm`", func(t *testing.T) {
		got := svc.buildTempCondition(76.0)
		assert.EqualValues(t, got, warm)
	})
}
func TestWeatherService_buildWindCondition(t *testing.T) {
	svc := &WeatherService{
		config:        &config.App{},
		WeatherClient: &openweather.Client{},
	}
	t.Run("Should return `calm`", func(t *testing.T) {
		got := svc.buildWindCondition(0)
		assert.EqualValues(t, got, calm)
	})
	t.Run("Should return `light air`", func(t *testing.T) {
		got := svc.buildWindCondition(1.3)
		assert.EqualValues(t, got, lightAir)
	})
	t.Run("Should return `light breeze`", func(t *testing.T) {
		got := svc.buildWindCondition(5)
		assert.EqualValues(t, got, lightBreeze)
	})
	t.Run("Should return `gentle breeze`", func(t *testing.T) {
		got := svc.buildWindCondition(9)
		assert.EqualValues(t, got, gentalBreeze)
	})
	t.Run("Should return `moderate breeze`", func(t *testing.T) {
		got := svc.buildWindCondition(13)
		assert.EqualValues(t, got, moderateBreeze)
	})
	t.Run("Should return `fresh breeze`", func(t *testing.T) {
		got := svc.buildWindCondition(19)
		assert.EqualValues(t, got, freshBreeze)
	})
	t.Run("Should return `strong breeze`", func(t *testing.T) {
		got := svc.buildWindCondition(25)
		assert.EqualValues(t, got, strongBreeze)
	})
	t.Run("Should return `near gale winds`", func(t *testing.T) {
		got := svc.buildWindCondition(33)
		assert.EqualValues(t, got, nearGale)
	})
	t.Run("Should return `gale winds`", func(t *testing.T) {
		got := svc.buildWindCondition(39)
		assert.EqualValues(t, got, gale)
	})
	t.Run("Should return `severe gale winds`", func(t *testing.T) {
		got := svc.buildWindCondition(47)
		assert.EqualValues(t, got, severeGale)
	})
	t.Run("Should return `storm winds`", func(t *testing.T) {
		got := svc.buildWindCondition(55)
		assert.EqualValues(t, got, storm)
	})
	t.Run("Should return `violent storm winds`", func(t *testing.T) {
		got := svc.buildWindCondition(64)
		assert.EqualValues(t, got, violentStorm)
	})
	t.Run("Should return `hurricane  winds`", func(t *testing.T) {
		got := svc.buildWindCondition(73)
		assert.EqualValues(t, got, hurricane)
	})
}
