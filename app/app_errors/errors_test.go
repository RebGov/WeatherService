package apperrors_test

import (
	"errors"
	"testing"
	apperrors "weathersvc/app/app_errors"

	"github.com/stretchr/testify/assert"
)

func TestErrors_CreateMissingConfigError(t *testing.T) {
	t.Run("", func(t *testing.T) {
		expected := errors.New("failed to start service: missing required config for `WEATHER_ID`")
		got := apperrors.CreateMissingConfigError("WEATHER_ID")
		assert.EqualError(t, got, expected.Error())
	})
}

func TestErrors_CreateInvalidRequestErrors(t *testing.T) {
	t.Run("", func(t *testing.T) {
		expected := errors.New("invalid request: latitude is out of range")
		got := apperrors.CreateInvalidRequestError("latitude is out of range")
		assert.EqualError(t, got, expected.Error())
	})
}
