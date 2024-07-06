package openweather

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"weathersvc/app/config"

	"github.com/stretchr/testify/assert"
)

type responseWriter struct {
	http.ResponseWriter
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.ResponseWriter.(http.Flusher).Flush() // Ensure the header is written before the body
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return 0, fmt.Errorf("simulated write error")
}

func Test_GetWeather(t *testing.T) {
	t.Run("Should fail to get weather for unvalid uri", func(t *testing.T) {
		conf := &config.App{
			Port: "fake",
			Env:  "testing",
			WeatherClientConfig: config.WeatherClientConfig{
				Host:  "fake",
				AppID: "fakefake",
			},
		}
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusBadRequest)
		}))
		defer testServer.Close()
		owmClient := NewClient(conf)
		resp, err := owmClient.GetWeather("0", "0")
		assert.EqualError(t, err, "error sending request: Get \"fake?appid=fakefake&lat=0&lon=0&units=imperial\": unsupported protocol scheme \"\"")
		assert.Nil(t, resp)
	})
	t.Run("Should return 401", func(t *testing.T) {
		conf := &config.App{
			Port: "fake",
			Env:  "testing",
			WeatherClientConfig: config.WeatherClientConfig{
				Host:  "",
				AppID: "fakefake",
			},
		}
		unauthorizedResponse := `{"cod": 401, "message": "Invalid API key."}`
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusUnauthorized)
			res.Header().Set("Content-Type", "application/json")
			res.Write([]byte(unauthorizedResponse))
		}))
		defer testServer.Close()
		conf.WeatherClientConfig.Host = testServer.URL
		owmClient := NewClient(conf)
		resp, err := owmClient.GetWeather("0", "0")
		assert.EqualError(t, err, "config `WEATHER_ID` is invalid")
		assert.Nil(t, resp)
	})
	t.Run("Should return 429", func(t *testing.T) {
		conf := &config.App{
			Port: "fake",
			Env:  "testing",
			WeatherClientConfig: config.WeatherClientConfig{
				Host:  "",
				AppID: "fakefake",
			},
		}
		toManyRequests := `{"cod": 429, "message": "Out of limit requsts."}`
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusTooManyRequests)
			res.Header().Set("Content-Type", "application/json")
			res.Write([]byte(toManyRequests))
		}))
		defer testServer.Close()
		conf.WeatherClientConfig.Host = testServer.URL
		owmClient := NewClient(conf)
		resp, err := owmClient.GetWeather("0", "0")
		assert.EqualError(t, err, "too many requests; limit reached")
		assert.Nil(t, resp)
	})
	t.Run("Should return 404", func(t *testing.T) {
		conf := &config.App{
			Port: "fake",
			Env:  "testing",
			WeatherClientConfig: config.WeatherClientConfig{
				Host:  "",
				AppID: "fakefake",
			},
		}
		requestResponse := `{"cod": 404, "message": "Coordinates not found."}`
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusNotFound)
			res.Header().Set("Content-Type", "application/json")
			res.Write([]byte(requestResponse))
		}))
		defer testServer.Close()
		conf.WeatherClientConfig.Host = testServer.URL
		owmClient := NewClient(conf)
		resp, err := owmClient.GetWeather("0", "0")
		assert.EqualError(t, err, "weather for coordinates not found")
		assert.Nil(t, resp)
	})
	t.Run("Should return 500", func(t *testing.T) {
		conf := &config.App{
			Port: "fake",
			Env:  "testing",
			WeatherClientConfig: config.WeatherClientConfig{
				Host:  "",
				AppID: "fakefake",
			},
		}
		requestResponse := `{"cod": 500, "message": "Some Internal Error."}`
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusInternalServerError)
			res.Header().Set("Content-Type", "application/json")
			res.Write([]byte(requestResponse))
		}))
		defer testServer.Close()
		conf.WeatherClientConfig.Host = testServer.URL
		owmClient := NewClient(conf)
		resp, err := owmClient.GetWeather("0", "0")
		assert.EqualError(t, err, "internal service error")
		assert.Nil(t, resp)
	})
	t.Run("Should return unable to unmarshal", func(t *testing.T) {
		conf := &config.App{
			Port: "fake",
			Env:  "testing",
			WeatherClientConfig: config.WeatherClientConfig{
				Host:  "",
				AppID: "fakefake",
			},
		}
		requestResponse := `{"cod": 500, "message": ome Internal Error."}`
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusInternalServerError)
			res.Header().Set("Content-Type", "application/json")
			res.Write([]byte(requestResponse))
		}))
		defer testServer.Close()
		conf.WeatherClientConfig.Host = testServer.URL
		owmClient := NewClient(conf)
		resp, err := owmClient.GetWeather("0", "0")
		assert.EqualError(t, err, "error unmarshalling response: invalid character 'o' looking for beginning of value")
		assert.Nil(t, resp)
	})
	t.Run("Should handle read error", func(t *testing.T) {
		conf := &config.App{
			Port: "fake",
			Env:  "testing",
			WeatherClientConfig: config.WeatherClientConfig{
				Host:  "",
				AppID: "fakefake",
			},
		}
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			// simulate a read response error
			rw := &responseWriter{res}
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			// Simulate a read error by writing to a custom response body (set to empty to avoid unmarshal errors)
			res.Write([]byte(`{}`))
			// utilize hiJacker to allow for customer read error
			hijacker, ok := res.(http.Hijacker)
			if !ok {
				t.Fatal("expected http.ResponseWriter to implement http.Hijacker")
			}
			conn, _, err := hijacker.Hijack()
			if err != nil {
				t.Fatalf("unexpected error hijacking the connection: %v", err)
			}
			conn.Close() // Closing the connection simulates a read error
		}))
		defer testServer.Close()
		conf.WeatherClientConfig.Host = testServer.URL
		owmClient := NewClient(conf)
		resp, err := owmClient.GetWeather("0", "0")
		assert.Nil(t, resp)
		assert.EqualError(t, err, "error reading response: unexpected EOF")
	})
	t.Run("Should return 200", func(t *testing.T) {
		conf := &config.App{
			Port: "fake",
			Env:  "testing",
			WeatherClientConfig: config.WeatherClientConfig{
				Host:  "",
				AppID: "fakefake",
			},
		}
		requestResponse := `{"coord":{"lon":0,"lat":0},"weather":[{"id":803,"main":"Clouds","description":"broken clouds","icon":"04d"}],"base":"stations","main":{"temp":75.83,"feels_like":76.78,"temp_min":75.83,"temp_max":75.83,"pressure":1013,"humidity":78,"sea_level":1013,"grnd_level":1013},"visibility":10000,"wind":{"speed":13.24,"deg":152,"gust":12.82},"clouds":{"all":54},"dt":1720197314,"sys":{"sunrise":1720159251,"sunset":1720202883},"timezone":0,"id":6295630,"name":"Globe","cod":200}`
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			res.Header().Set("Content-Type", "application/json")
			res.Write([]byte(requestResponse))
		}))
		defer testServer.Close()
		conf.WeatherClientConfig.Host = testServer.URL
		owmClient := NewClient(conf)
		resp, err := owmClient.GetWeather("0", "0")
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func Test_ApiTest(t *testing.T) {
	t.Run("Should return 401", func(t *testing.T) {
		conf := &config.App{
			Port: "fake",
			Env:  "testing",
			WeatherClientConfig: config.WeatherClientConfig{
				Host:  "",
				AppID: "fakefake",
			},
		}
		unauthorizedResponse := `{"cod": 401, "message": "Invalid API key."}`
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusUnauthorized)
			res.Header().Set("Content-Type", "application/json")
			res.Write([]byte(unauthorizedResponse))
		}))
		defer testServer.Close()
		conf.WeatherClientConfig.Host = testServer.URL
		owmClient := NewClient(conf)
		err := owmClient.ApiTest()
		assert.EqualError(t, err, "config `WEATHER_ID` is invalid")
	})

	t.Run("Should return 200", func(t *testing.T) {
		conf := &config.App{
			Port: "fake",
			Env:  "testing",
			WeatherClientConfig: config.WeatherClientConfig{
				Host:  "",
				AppID: "fakefake",
			},
		}
		requestResponse := `{"coord":{"lon":0,"lat":0},"weather":[{"id":803,"main":"Clouds","description":"broken clouds","icon":"04d"}],"base":"stations","main":{"temp":75.83,"feels_like":76.78,"temp_min":75.83,"temp_max":75.83,"pressure":1013,"humidity":78,"sea_level":1013,"grnd_level":1013},"visibility":10000,"wind":{"speed":13.24,"deg":152,"gust":12.82},"clouds":{"all":54},"dt":1720197314,"sys":{"sunrise":1720159251,"sunset":1720202883},"timezone":0,"id":6295630,"name":"Globe","cod":200}`
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			res.Header().Set("Content-Type", "application/json")
			res.Write([]byte(requestResponse))
		}))
		defer testServer.Close()
		conf.WeatherClientConfig.Host = testServer.URL
		owmClient := NewClient(conf)
		err := owmClient.ApiTest()
		assert.NoError(t, err)
	})

}
