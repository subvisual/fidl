package bank

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus/pkg/crypto"
	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/http/jsend"
)

func (s *Server) JSON(w http.ResponseWriter, r *http.Request, code int, value any) {
	var status int
	var body any

	if err, ok := value.(error); ok {
		switch {
		case errors.Is(err, ErrInsufficientFunds):
			status, body = http.StatusForbidden, envelope{"bank": "insufficient funds"}
		case errors.Is(err, ErrLockedFunds):
			status, body = http.StatusUnauthorized, envelope{"bank": "locked funds"}
		case errors.Is(err, ErrTransactionNotAllowed):
			status, body = http.StatusUnauthorized, envelope{"bank": "transaction not allowed"}
		default:
			s.HTTP.JSON(w, r, code, value)
			return
		}

		s.HTTP.LogDebug(r, err)
	} else {
		s.HTTP.JSON(w, r, code, value)
		return
	}

	fidl.SetHeaders(w, status)
	payload := jsend.Fail(body)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		s.HTTP.LogError(r, err)
	}
}

func ParseHeader(r *http.Request) (*crypto.Signature, address.Address, string, error) {
	dataSig := []byte(r.Header.Get("sig"))
	dataPub := []byte(r.Header.Get("pub"))
	dataMsg := []byte(r.Header.Get("msg"))

	var str string
	if err := json.Unmarshal(dataSig, &str); err != nil {
		return nil, address.Address{}, "", fmt.Errorf("failed to parse signature from header: %w", err)
	}

	var sig crypto.Signature
	if err := json.Unmarshal([]byte(str), &sig); err != nil {
		return nil, address.Address{}, "", fmt.Errorf("failed to parse signature from header: %w", err)
	}

	var pub string
	if err := json.Unmarshal(dataPub, &pub); err != nil {
		return nil, address.Address{}, "", fmt.Errorf("failed to parse public key from header: %w", err)
	}

	addr, err := address.NewFromString(pub)
	if err != nil {
		return nil, address.Address{}, "", fmt.Errorf("failed to parse address from header: %w", err)
	}

	var msg string
	if err := json.Unmarshal(dataMsg, &msg); err != nil {
		return nil, address.Address{}, "", fmt.Errorf("failed to parse message from header: %w", err)
	}

	return &sig, addr, msg, nil
}
