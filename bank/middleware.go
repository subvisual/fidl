package bank

import (
	"context"
	"net/http"

	"github.com/subvisual/fidl/crypto"
)

type ctxKey int

const (
	CtxKeyAddress ctxKey = iota
)

func AuthenticationCtx() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sig, addr, msg, err := ParseHeader(r)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			if err := crypto.Verify(sig, *addr.Address, msg); err != nil {
				http.Error(w, "failed to verify signature", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, CtxKeyAddress, addr)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
