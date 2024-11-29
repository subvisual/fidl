package bank

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/subvisual/fidl/http/jsend"
	"github.com/subvisual/fidl/types"
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
		case errors.Is(err, ErrOperationNotAllowed):
			status, body = http.StatusUnauthorized, envelope{"bank": "operation not allowed"}
		case errors.Is(err, ErrNothingToRefund):
			status, body = http.StatusUnprocessableEntity, envelope{"bank": "nothing to refund"}
		case errors.Is(err, ErrAuthNotFound):
			status, body = http.StatusNotFound, envelope{"bank": "no valid authorization"}
		case errors.Is(err, ErrAuthLocked):
			status, body = http.StatusNotFound, envelope{"bank": "authorization is locked"}
		default:
			s.Server.JSON(w, r, code, value)
			return
		}

		s.LogDebug(r, err)
	} else {
		s.Server.JSON(w, r, code, value)
		return
	}

	SetHeaders(w, status)
	payload := jsend.Fail(body)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		s.LogError(r, err)
	}
}

func ParseHeader(r *http.Request) (*types.Signature, types.Address, []byte, error) {
	dataSig := r.Header.Get("sig")
	dataPub := r.Header.Get("pub")
	dataMsg := r.Header.Get("msg")

	binSig, err := hex.DecodeString(dataSig)
	if err != nil {
		return nil, types.Address{}, nil, fmt.Errorf("failed to decode signature string: %w", err)
	}

	binMsg, err := hex.DecodeString(dataMsg)
	if err != nil {
		return nil, types.Address{}, nil, fmt.Errorf("failed to decode message string: %w", err)
	}

	var sig types.Signature
	if err = sig.UnmarshalBinary(binSig); err != nil {
		return nil, types.Address{}, nil, fmt.Errorf("failed to unmarshal binary signature: %w", err)
	}

	addr, err := types.NewAddressFromString(dataPub)
	if err != nil {
		return nil, types.Address{}, nil, fmt.Errorf("failed to parse address from header: %w", err)
	}

	return &sig, addr, binMsg, nil
}
