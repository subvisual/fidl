package validation

import (
	"errors"

	vtypes "github.com/filecoin-project/venus/venus-shared/actors/types"
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

func IsValidAddress(fl validator.FieldLevel) bool {
	_, err := vtypes.ParseEthAddress(fl.Field().String())
	if err == nil {
		return true
	}

	filAddr, err := types.NewAddressFromString(fl.Field().String())
	if err != nil {
		return false
	}

	_, err = vtypes.EthAddressFromFilecoinAddress(*filAddr.Address)

	return err == nil
}

func IsValidFIL(fl validator.FieldLevel) bool {
	if fil, ok := fl.Field().Interface().(types.FIL); !ok || !(fil.Int.Sign() == 1) {
		return false
	}

	return true
}
