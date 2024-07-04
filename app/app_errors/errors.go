package apperrors

import (
	"errors"
	"fmt"
)

//nolint:gochecknoglobals // 20240702BG allow
var (
	ErrMissingConfig        = errors.New("failed to start service: missing required config")
	ErrInvalidRequest       = errors.New("latitude and longitude cannot be set to zero values")
	ErrInvalidOWMAppID      = errors.New("config `WEATHER_ID` is invalid")
	ErrInternalServiceError = errors.New("internal service error")
	ErrTooManyRequests      = errors.New("too many requests; limit reached")
	ErrNotFound             = errors.New("weather for coordinates not found")
)

func CreateMissingConfigError(e error, v string) error {
	return fmt.Errorf("%s for `%s`", e.Error(), v)
}
