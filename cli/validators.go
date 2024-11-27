package cli

import (
	"fmt"

	"github.com/subvisual/fidl/validation"
)

func RegisterValidators(cli CLI) error {
	if err := cli.Validate.RegisterValidation("is-filecoin-address", validation.IsFilecoinAddress); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
