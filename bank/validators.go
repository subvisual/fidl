package bank

import (
	"github.com/subvisual/fidl/validation"
	"go.uber.org/zap"
)

func (s *Server) RegisterValidators() {
	// Register validators
	if err := s.Validate.RegisterValidation("is-filecoin-address", validation.IsFilecoinAddress); err != nil {
		s.Log.Fatal("Unable to register is-filecoin-address validator", zap.String("name", "is-filecoin-address"))
	}

	if err := s.Validate.RegisterValidation("is-valid-fil", validation.IsValidFIL); err != nil {
		s.Log.Fatal("Unable to register is-valid-fil validator", zap.String("name", "is-valid-fil"))
	}
}
