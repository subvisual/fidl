package proxy

import (
	"github.com/google/uuid"
	"github.com/subvisual/fidl/http"
	"github.com/subvisual/fidl/types"
)

type Server struct {
	*http.Server
	Bank          map[string]Bank
	ExternalRoute Route
	Forwarder     *Forwarder
	Provider      Provider
	Wallet        types.Wallet
}

type RetrievalParams struct {
	Authorization uuid.UUID `validate:"required"`
	Bank          string    `validate:"required,url"`
}

type BankListResponse struct {
	Cost types.FIL `json:"cost"`
	URL  string    `json:"url"`
}
