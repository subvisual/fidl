package cli

import (
	"github.com/go-playground/validator/v10"
)

type CLI struct {
	Validate *validator.Validate
}

func NewCLI(validate *validator.Validate) CLI {
	return CLI{
		Validate: validate,
	}
}
