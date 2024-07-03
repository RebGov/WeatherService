package apperrors

import (
	"errors"
	"fmt"
)

//nolint:gochecknoglobals // 20240702BG allow
var (
	ErrMissingConfig        = errors.New("failed to start service: missing required config")
	ErrInvalidRequest       = errors.New("latitude and longitude cannot be set to zero values")
	ErrInvalidOWMAppID      = errors.New("config `WEATHER_APPID` is invalid")
	ErrInternalServiceError = errors.New("internal service error")
)

func CreateMissingConfigError(e error, v string) error {
	return fmt.Errorf("%s for `%s`", e.Error(), v)
}