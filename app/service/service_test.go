package service

import (
	"context"
	"testing"
	apperrors "weathersvc/app/app_errors"
	"weathersvc/app/config"
	"weathersvc/app/models"
	ownMock "weathersvc/mocks/open_weather"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_NewService(t *testing.T) {
	ctx := context.Background()
	conf := &config.App{
		Port: "8081",
		Env:  "testing",
		WeatherClientConfig: config.WeatherClientConfig{
			Host:  "",
			AppID: "fakeID",
		},
		ServiceURL: "fakevalue",
	}
	t.Run("Should not fail to create new service", func(t *testing.T) {
		got := NewService(ctx, conf)
		assert.NotNil(t, got)
	})
}
func TestService_buildTempCondition(t *testing.T) {
	svc := service{
		Config:        &config.App{},
		WeatherClient: &ownMock.MockClient{},
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
func TestService_buildWindCondition(t *testing.T) {
	svc := service{
		Config:        &config.App{},
		WeatherClient: &ownMock.MockClient{},
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
func TestSevice_ValidateSvc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	owm := ownMock.NewMockClient(ctrl)
	svc := service{
		Config: &config.App{
			Port: "item",
			Env:  "testing",
			WeatherClientConfig: config.WeatherClientConfig{
				Host:  "item",
				AppID: "item",
			},
			ServiceURL: "item",
		},
		WeatherClient: owm,
	}
	t.Run("Should pass validation", func(t *testing.T) {
		owm.EXPECT().ApiTest().Return(nil)
		err := svc.ValidateSvc(context.Background())
		assert.NoError(t, err)
	})
	t.Run("Should fail validation when user started svc with invalid weather-appid", func(t *testing.T) {
		owm.EXPECT().ApiTest().Return(apperrors.ErrInvalidOWMAppID)
		err := svc.ValidateSvc(context.Background())
		assert.EqualError(t, err, apperrors.ErrInvalidOWMAppID.Error())
	})
}

func TestService_GetWeather(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	owm := ownMock.NewMockClient(ctrl)
	svc := service{
		Config: &config.App{
			Port: "item",
			Env:  "testing",
			WeatherClientConfig: config.WeatherClientConfig{
				Host:  "item",
				AppID: "item",
			},
			ServiceURL: "item",
		},
		WeatherClient: owm,
	}
	t.Run("Should return valid response for valid lat/lon", func(t *testing.T) {
		expectResp := WeatherCond{
			Temp:      hot,
			Condition: "few clouds",
			Wind:      calm,
		}
		desc1 := models.Weather{
			Description: "few clouds",
		}
		desc := []models.Weather{}
		desc = append(desc, desc1)
		owm.EXPECT().GetWeather(gomock.Any(), gomock.Any()).Return(&models.WeatherResponse{
			Weather: desc,
			Main: models.Main{
				FeelsLike: 90.4,
			},
			Wind: models.Wind{
				Speed: 0,
			},
			Cod: 200,
		}, nil)
		got, gErr := svc.GetWeather(context.Background(), 0, 0)
		assert.NoError(t, gErr)
		assert.EqualValues(t, got.Temp, expectResp.Temp)
		assert.EqualValues(t, got.Condition, expectResp.Condition)
		assert.EqualValues(t, got.Wind, expectResp.Wind)
	})
	t.Run("Should return err 429", func(t *testing.T) {
		expectResp := WeatherCond{
			Temp:      "",
			Condition: "",
			Wind:      "",
		}
		owm.EXPECT().GetWeather(gomock.Any(), gomock.Any()).Return(&models.WeatherResponse{}, apperrors.ErrTooManyRequests)
		got, gErr := svc.GetWeather(context.Background(), 0, 0)
		assert.EqualError(t, gErr, "too many requests; limit reached")
		assert.EqualValues(t, got.Temp, expectResp.Temp)
		assert.EqualValues(t, got.Condition, expectResp.Condition)
		assert.EqualValues(t, got.Wind, expectResp.Wind)
	})
	t.Run("Should not return err when condition is missing", func(t *testing.T) {
		expectResp := WeatherCond{
			Temp:      hot,
			Condition: "unknown",
			Wind:      calm,
		}

		owm.EXPECT().GetWeather(gomock.Any(), gomock.Any()).Return(&models.WeatherResponse{
			Main: models.Main{
				FeelsLike: 90.4,
			},
			Wind: models.Wind{
				Speed: 0,
			},
			Cod: 200,
		}, nil)
		got, gErr := svc.GetWeather(context.Background(), 0, 0)
		assert.NoError(t, gErr)
		assert.EqualValues(t, got.Temp, expectResp.Temp)
		assert.EqualValues(t, got.Condition, expectResp.Condition)
		assert.EqualValues(t, got.Wind, expectResp.Wind)
	})

}
