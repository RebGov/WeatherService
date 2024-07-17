package openweather

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	apperrors "weathersvc/app/app_errors"
	"weathersvc/app/config"
	"weathersvc/app/models"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
)

type Client interface {
	ApiTest() error
	GetWeather(lat, lon string) (*models.WeatherResponse, error)
}
type client struct {
	client *http.Client
	host   string
	appId  string
}

func NewClient(conf *config.App) Client {
	httpClient := &http.Client{
		Transport: &nethttp.Transport{},
	}
	return &client{
		client: httpClient,
		host:   conf.WeatherClientConfig.Host,
		appId:  conf.WeatherClientConfig.AppID,
	}
}

func (c *client) ApiTest() error {
	_, err := c.GetWeather("0", "0")
	if err != nil {
		return err
	}
	return nil
}

// GetWeather takes in string: latitude longitude
func (c *client) GetWeather(lat, long string) (*models.WeatherResponse, error) {
	u, err := url.Parse(c.host)
	if err != nil {
		return nil, nil
	}
	query := url.Values{}
	query.Add("lat", lat)
	query.Add("lon", long)
	query.Add("units", "imperial")
	query.Add("appid", c.appId)
	u.RawQuery = query.Encode()
	// Create the HTTP request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	// Set headers if necessary
	req.Header.Set("Content-Type", "application/json")
	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	var data *models.WeatherResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("error on Unmarshall: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}
	switch data.Cod {
	case 200:
		return data, nil
	case 401:
		return nil, apperrors.ErrInvalidOWMAppID
	case 404:
		return nil, apperrors.ErrNotFound
	case 429:
		return nil, apperrors.ErrTooManyRequests
	default:
		return nil, apperrors.ErrInternalServiceError
	}
}
