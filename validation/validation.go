package validation

import (
	"errors"

	"go.uber.org/zap"
)

var (
	ErrInvalidContentLength = errors.New("invalid content length")
	ErrInvalidMimeType      = errors.New("invalid mime type")
)

type StoreValidator struct {
	log *zap.Logger
}

func New(log *zap.Logger) *StoreValidator {
	return &StoreValidator{log: log}
}
