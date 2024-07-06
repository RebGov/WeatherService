package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	apperrors "weathersvc/app/app_errors"
	"weathersvc/app/config"
	"weathersvc/app/service"
	mock_service "weathersvc/mocks/service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServer_NewServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := mock_service.NewMockService(ctrl)
	conf := &config.App{
		Port: "8081",
		Env:  "testing",
		WeatherClientConfig: config.WeatherClientConfig{
			Host:  "",
			AppID: "fakeID",
		},
		ServiceURL: "fakevalue",
	}
	t.Run("Should build server", func(t *testing.T) {
		got := NewServer(conf, svc)
		assert.NotNil(t, got)
	})
}

func TestGetWeatherHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockService(ctrl)
	t.Run("Should pass 200", func(t *testing.T) {
		mockService.EXPECT().GetWeather(gomock.Any(), gomock.Any(), gomock.Any()).Return(service.WeatherCond{
			Temp:      "hot",
			Condition: "few clouds",
			Wind:      "calm",
		}, nil)
		reqBody := DecimalRequest{
			Latitude:  -1,
			Longitude: 1,
		}
		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)
		req := httptest.NewRequest("GET", "/weather/get/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(weatherHandler(mockService))
		handler.ServeHTTP(rr, req)
		var respBody Response
		err = json.NewDecoder(rr.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Equal(t, "hot", respBody.Temp)
		assert.Equal(t, "few clouds", respBody.Condition)
		assert.Equal(t, "calm", respBody.Wind)
	})
	t.Run("Should fail 500", func(t *testing.T) {
		mockService.EXPECT().GetWeather(gomock.Any(), gomock.Any(), gomock.Any()).Return(service.WeatherCond{
			Temp:      "",
			Condition: "",
			Wind:      "",
		}, apperrors.ErrInternalServiceError)
		reqBody := DecimalRequest{
			Latitude:  1,
			Longitude: 1,
		}
		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)
		req := httptest.NewRequest("GET", "/weather/get/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(weatherHandler(mockService))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "internal service error")
	})
	t.Run("Should fail 429", func(t *testing.T) {
		mockService.EXPECT().GetWeather(gomock.Any(), gomock.Any(), gomock.Any()).Return(service.WeatherCond{
			Temp:      "",
			Condition: "",
			Wind:      "",
		}, apperrors.ErrTooManyRequests)
		reqBody := DecimalRequest{
			Latitude:  1,
			Longitude: 1,
		}
		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)
		req := httptest.NewRequest("GET", "/weather/get/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(weatherHandler(mockService))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		assert.Contains(t, rr.Body.String(), "too many requests; limit reached")
	})
	t.Run("Should fail 404", func(t *testing.T) {
		mockService.EXPECT().GetWeather(gomock.Any(), gomock.Any(), gomock.Any()).Return(service.WeatherCond{
			Temp:      "",
			Condition: "",
			Wind:      "",
		}, apperrors.ErrNotFound)
		reqBody := DecimalRequest{
			Latitude:  -1,
			Longitude: 1,
		}
		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)
		req := httptest.NewRequest("GET", "/weather/get/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(weatherHandler(mockService))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "coordinates not found")
	})

	t.Run("Should fail 400 for lat and lon set to zero", func(t *testing.T) {
		reqBody := DecimalRequest{
			Latitude:  0,
			Longitude: 0,
		}
		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)
		req := httptest.NewRequest("GET", "/weather/get/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(weatherHandler(mockService))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid request: latitude and longitude missing or null")
	})
	t.Run("Should fail 400 for lat out of range", func(t *testing.T) {
		reqBody := DecimalRequest{
			Latitude:  -90.01,
			Longitude: 0,
		}
		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)
		req := httptest.NewRequest("GET", "/weather/get/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(weatherHandler(mockService))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid request: latitude is out of range")
	})
	t.Run("Should fail 400 for lon out of range", func(t *testing.T) {
		reqBody := DecimalRequest{
			Latitude:  90,
			Longitude: -180.01,
		}
		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)
		req := httptest.NewRequest("GET", "/weather/get/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(weatherHandler(mockService))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid request: longitude is out of range")
	})
	t.Run("Should fail 400 missing body", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/weather/get/", bytes.NewBuffer([]byte{}))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(weatherHandler(mockService))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "request body missing: see `https://github.com/RebGov/WeatherService")
	})
}

func TestServer_Open(t *testing.T) {
	conf := &config.App{Port: "0"} // Use port "0" to let the system choose an available port
	s := NewServer(conf, nil).(*server)
	t.Run("Should not be listing on port 0", func(t *testing.T) {
		done := make(chan error)
		go func() {
			done <- s.Open()
		}()
		select {
		case <-time.After(1 * time.Second): // Wait up to 1 second for the server to start
		case err := <-done:
			t.Fatalf("Failed to open server: %v", err)
		}
		// Check if the server is listening on a port
		port := s.Port()
		assert.NotEqual(t, 0, port, "Server should not be listening on port 0")
	})
}

func TestServer_Close(t *testing.T) {
	conf := &config.App{Port: "0"} // Use port "0" to let the system choose an available port
	s := NewServer(conf, nil).(*server)
	// Use a goroutine to open the server as it will block
	go func() {
		err := s.Open()
		assert.NoError(t, err)
	}()
	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)
	// Close the server and check for errors
	err := s.Close()
	assert.NoError(t, err)
}

func TestServer_Port(t *testing.T) {
	conf := &config.App{Port: "0"} // Use port "0" to let the system choose an available port
	s := NewServer(conf, nil).(*server)

	// Check the port before the server is opened
	assert.Equal(t, 0, s.Port())
	// Use a goroutine to open the server as it will block
	go func() {
		err := s.Open()
		assert.NoError(t, err)
	}()
	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)
	// Check if the server is listening on a port
	assert.NotEqual(t, 0, s.Port())
}
