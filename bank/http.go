package bank

import (
	"encoding/json"
	"errors"
	"net/http"

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
