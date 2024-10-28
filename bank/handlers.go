package bank

import (
	"net/http"

	"github.com/filecoin-project/go-address"
	"github.com/go-chi/chi/v5"
)

type envelope map[string]any

func (s *Server) Routes(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Use(AuthenticationCtx())
		r.Post("/register", s.handleRegisterProxy)
		r.Post("/deposit", s.handleDeposit)
		r.Post("/withdraw", s.handleWithdraw)
		r.Get("/balance", s.handleBalance)
		r.Post("/authorize", s.handleAuthorize)
		r.Post("/redeem", s.handleRedeem)
	})
}

func (s *Server) handleRegisterProxy(w http.ResponseWriter, r *http.Request) {
	var params RegisterParams

	address, ok := r.Context().Value(CtxKeyAddress).(address.Address)
	if !ok {
		s.JSON(w, r, http.StatusBadRequest, "failed to parse header signature")
		return
	}

	if err := s.HTTP.DecodeJSON(w, r, &params); err != nil {
		s.JSON(w, r, http.StatusBadRequest, envelope{"message": err.Error()})
		return
	}

	if err := s.HTTP.Validate.Struct(params); err != nil {
		s.JSON(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	if err := s.BankService.RegisterProxy(params.ID, address.String(), params.Price); err != nil {
		s.JSON(w, r, http.StatusInternalServerError, err)
		return
	}

	s.JSON(w, r, http.StatusOK, envelope{"bank": "proxy registered"})
}

func (s *Server) handleDeposit(w http.ResponseWriter, r *http.Request) {
	var params DepositParams

	address, ok := r.Context().Value(CtxKeyAddress).(address.Address)
	if !ok {
		s.JSON(w, r, http.StatusBadRequest, "failed to parse header signature")
		return
	}

	if err := s.HTTP.DecodeJSON(w, r, &params); err != nil {
		s.JSON(w, r, http.StatusBadRequest, envelope{"message": err.Error()})
		return
	}

	if err := s.HTTP.Validate.Struct(params); err != nil {
		s.JSON(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	fil, err := s.BankService.Deposit(address.String(), params.Amount)
	if err != nil {
		s.JSON(w, r, http.StatusInternalServerError, err)
		return
	}

	s.JSON(w, r, http.StatusOK, envelope{"fil": fil})
}

func (s *Server) handleWithdraw(w http.ResponseWriter, r *http.Request) {
	var params WithdrawParams

	address, ok := r.Context().Value(CtxKeyAddress).(address.Address)
	if !ok {
		s.JSON(w, r, http.StatusBadRequest, "failed to parse header signature")
		return
	}

	if err := s.HTTP.DecodeJSON(w, r, &params); err != nil {
		s.JSON(w, r, http.StatusBadRequest, envelope{"message": err.Error()})
		return
	}

	if err := s.HTTP.Validate.Struct(params); err != nil {
		s.JSON(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	fil, err := s.BankService.Withdraw(address.String(), params.Destination, params.Amount)
	if err != nil {
		s.JSON(w, r, http.StatusInternalServerError, err)
		return
	}

	s.JSON(w, r, http.StatusOK, envelope{"fil": fil})
}

func (s *Server) handleBalance(w http.ResponseWriter, r *http.Request) {
	address, ok := r.Context().Value(CtxKeyAddress).(address.Address)
	if !ok {
		s.JSON(w, r, http.StatusBadRequest, "failed to parse header signature")
		return
	}

	fil, err := s.BankService.Balance(address.String())
	if err != nil {
		s.JSON(w, r, http.StatusInternalServerError, err)
		return
	}

	s.JSON(w, r, http.StatusOK, envelope{"fil": fil})
}

func (s *Server) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	var params AuthorizeParams

	_, ok := r.Context().Value(CtxKeyAddress).(address.Address)
	if !ok {
		s.JSON(w, r, http.StatusBadRequest, "failed to parse header signature")
		return
	}

	if err := s.HTTP.DecodeJSON(w, r, &params); err != nil {
		s.JSON(w, r, http.StatusBadRequest, envelope{"message": err.Error()})
		return
	}

	if err := s.HTTP.Validate.Struct(params); err != nil {
		s.JSON(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	/* TODO */

	s.JSON(w, r, http.StatusOK, envelope{"bank": "TODO: authorize"})
}

func (s *Server) handleRedeem(w http.ResponseWriter, r *http.Request) {
	var params RedeemParams

	_, ok := r.Context().Value(CtxKeyAddress).(address.Address)
	if !ok {
		s.JSON(w, r, http.StatusBadRequest, "failed to parse header signature")
		return
	}

	if err := s.HTTP.DecodeJSON(w, r, &params); err != nil {
		s.JSON(w, r, http.StatusBadRequest, envelope{"message": err.Error()})
		return
	}

	if err := s.HTTP.Validate.Struct(params); err != nil {
		s.JSON(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	/* TODO */

	s.JSON(w, r, http.StatusOK, envelope{"bank": "TODO: redeem"})
}
