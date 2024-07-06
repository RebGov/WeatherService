package apperrors

import (
	"errors"
	"fmt"
)

//nolint:gochecknoglobals // 20240702BG allow
var (
	ErrMissingConfig        = errors.New("failed to start service: missing required config")
	ErrInvalidRequest       = errors.New("invalid request")
	ErrInvalidOWMAppID      = errors.New("config `WEATHER_ID` is invalid")
	ErrInternalServiceError = errors.New("internal service error")
	ErrTooManyRequests      = errors.New("too many requests; limit reached")
	ErrNotFound             = errors.New("weather for coordinates not found")
	ErrNoBody               = errors.New("request body missing: see `https://github.com/RebGov/WeatherService`")
)

// CreateMissingConfigError combines the missing environment config error and reason
func CreateMissingConfigError(v string) error {
	return fmt.Errorf("%s for `%s`", ErrMissingConfig.Error(), v)
}

// CreateInvalidRequestError combines the invalid request error and reason
func CreateInvalidRequestError(v string) error {
	return fmt.Errorf("%s: %s", ErrInvalidRequest.Error(), v)
}
