package bank

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus/pkg/crypto"
	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/http/jsend"
)

func SetHeaders(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
}

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

	SetHeaders(w, status)
	payload := jsend.Fail(body)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		s.HTTP.LogError(r, err)
	}
}

func ParseHeader(r *http.Request) (*crypto.Signature, fidl.Address, []byte, error) {
	dataSig := r.Header.Get("sig")
	dataPub := r.Header.Get("pub")
	dataMsg := r.Header.Get("msg")

	binSig, err := hex.DecodeString(dataSig)
	if err != nil {
		return nil, fidl.Address{}, nil, fmt.Errorf("failed to decode signature string: %w", err)
	}

	binMsg, err := hex.DecodeString(dataMsg)
	if err != nil {
		return nil, fidl.Address{}, nil, fmt.Errorf("failed to decode message string: %w", err)
	}

	var sig crypto.Signature
	if err = sig.UnmarshalBinary(binSig); err != nil {
		return nil, fidl.Address{}, nil, fmt.Errorf("failed to unmarshal binary signature: %w", err)
	}

	addr, err := address.NewFromString(dataPub)
	if err != nil {
		return nil, fidl.Address{}, nil, fmt.Errorf("failed to parse address from header: %w", err)
	}

	return &sig, fidl.Address{Address: &addr}, binMsg, nil
}
