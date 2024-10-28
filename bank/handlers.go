package bank

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/subvisual/fidl"
)

type envelope map[string]any

func (s *Server) Routes(r chi.Router) {
	r.Post("/register", s.handleRegisterProxy)
	r.Post("/deposit", func(w http.ResponseWriter, r *http.Request) {
		s.handleTransaction(w, r, s.BankService.Deposit)
	})
	r.Post("/withdraw", func(w http.ResponseWriter, r *http.Request) {
		s.handleTransaction(w, r, s.BankService.Withdraw)
	})
	r.Get("/balance", s.handleBalance)
	r.Post("/authorize", s.handleAuthorize)
	r.Post("/redeem", s.handleRedeem)
}

func (s *Server) handleRegisterProxy(w http.ResponseWriter, r *http.Request) {
	var params RegisterProxyParams

	if err := s.HTTP.DecodeJSON(w, r, &params); err != nil {
		s.JSON(w, r, http.StatusBadRequest, envelope{"message": err.Error()})
		return
	}

	if err := s.HTTP.Validate.Struct(params); err != nil {
		s.JSON(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	/*
		TODO: verify that the signature matches the pubkey
	*/

	err := s.BankService.RegisterProxy(params.SpID, params.PublicKey, params.Price)
	if err != nil {
		s.JSON(w, r, http.StatusInternalServerError, err)
		return
	}

	s.JSON(w, r, http.StatusOK, envelope{"bank": "proxy registered"})
}

func (s *Server) handleTransaction(w http.ResponseWriter, r *http.Request, transactionFn func(address string, amount fidl.FIL) (fidl.FIL, error)) {
	var params TransactionParams

	if err := s.HTTP.DecodeJSON(w, r, &params); err != nil {
		s.JSON(w, r, http.StatusBadRequest, envelope{"message": err.Error()})
		return
	}

	if err := s.HTTP.Validate.Struct(params); err != nil {
		s.JSON(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	/*
		TODO: verify that the signature matches the pubkey
	*/

	fil, err := transactionFn(params.PublicKey, params.Amount)
	if err != nil {
		s.JSON(w, r, http.StatusInternalServerError, err)
		return
	}

	s.JSON(w, r, http.StatusOK, envelope{"fil": fil})
}

func (s *Server) handleBalance(w http.ResponseWriter, r *http.Request) {
	var params BalanceParams

	if err := s.HTTP.DecodeJSON(w, r, &params); err != nil {
		s.JSON(w, r, http.StatusBadRequest, envelope{"message": err.Error()})
		return
	}

	if err := s.HTTP.Validate.Struct(params); err != nil {
		s.JSON(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	/*
		TODO: verify that the signature matches the pubkey
	*/

	fil, err := s.BankService.Balance(params.PublicKey)
	if err != nil {
		s.JSON(w, r, http.StatusInternalServerError, err)
		return
	}

	s.JSON(w, r, http.StatusOK, envelope{"fil": fil})
}

func (s *Server) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	s.JSON(w, r, http.StatusOK, envelope{"bank": "TODO: authorize"})
}

func (s *Server) handleRedeem(w http.ResponseWriter, r *http.Request) {
	s.JSON(w, r, http.StatusOK, envelope{"bank": "TODO: redeem"})
}
