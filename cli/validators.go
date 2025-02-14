package cli

import (
	"fmt"

	"github.com/subvisual/fidl/validation"
)

func (cli *CLI) RegisterValidators() error {
	if err := cli.Validate.RegisterValidation("is-filecoin-address", validation.IsFilecoinAddress); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := cli.Validate.RegisterValidation("is-valid-address", validation.IsValidAddress); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
