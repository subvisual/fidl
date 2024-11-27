package validation

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/subvisual/fidl/types"
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

func IsFilecoinAddress(fl validator.FieldLevel) bool {
	if _, err := types.NewAddressFromString(fl.Field().String()); err != nil {
		return false
	}

	return true
}

func IsValidFIL(fl validator.FieldLevel) bool {
	if fil, ok := fl.Field().Interface().(types.FIL); !ok || !(fil.Int.Sign() == 1) {
		return false
	}

	return true
}
