package bank

import (
	"github.com/subvisual/fidl/http"
)

type Server struct {
	*http.Server

	BankService Service
}

type Service interface {
}
