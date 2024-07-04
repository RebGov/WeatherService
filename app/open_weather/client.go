package openweather

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	apperrors "weathersvc/app/app_errors"
	"weathersvc/app/config"
	"weathersvc/app/models"
)

type Client struct {
	client *http.Client
	host   string
	appId  string
}

func NewClient(conf *config.App) *Client {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12, // Set minimum TLS version to TLS 1.2
		MaxVersion: tls.VersionTLS13, // Set maximum TLS version to TLS 1.3
	}
	// Create an HTTP client with the TLS configuration
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{
		Transport: transport,
	}
	return &Client{
		client: client,
		host:   conf.WeatherClientConfig.Host,
		appId:  conf.WeatherClientConfig.AppID,
	}
}

func (c *Client) ApiTest() error {
	_, err := c.GetWeather("0", "0")
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetWeather(lat, long string) (*models.WeatherResponse, error) {
	u, err := url.Parse(c.host)
	if err != nil {
		return nil, nil
	}
	query := url.Values{}
	query.Add("lat", lat)
	query.Add("lon", long)
	query.Add("appid", c.appId)
	u.RawQuery = query.Encode()
	// Create the HTTP request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	// Set headers if necessary
	req.Header.Set("Content-Type", "application/json")
	// log.Printf("Request: %+v", req)
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	// log.Printf("Body: `%s`", string(body))
	var data *models.WeatherResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("error on Unmarshall: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}
	log.Printf("unmarshal: %+v", data)
	switch data.Cod {
	case 200:
		return data, nil
	case 401:
		return nil, apperrors.ErrInvalidOWMAppID
	default:
		return nil, apperrors.ErrInternalServiceError
	}
}
