package bank

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/subvisual/fidl/types"
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
		r.Get("/refund", s.handleRefund)
		r.Post("/redeem", s.handleRedeem)
		r.Post("/verify", s.handleVerify)
	})
}

func (s *Server) handleRegisterProxy(w http.ResponseWriter, r *http.Request) {
	var params RegisterParams

	address, ok := r.Context().Value(CtxKeyAddress).(types.Address)
	if !ok {
		s.JSON(w, r, http.StatusBadRequest, "failed to parse header address")
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

	address, ok := r.Context().Value(CtxKeyAddress).(types.Address)
	if !ok {
		s.JSON(w, r, http.StatusBadRequest, "failed to parse header address")
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

	address, ok := r.Context().Value(CtxKeyAddress).(types.Address)
	if !ok {
		s.JSON(w, r, http.StatusBadRequest, "failed to parse header address")
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
	address, ok := r.Context().Value(CtxKeyAddress).(types.Address)
	if !ok {
		s.JSON(w, r, http.StatusBadRequest, "failed to parse header address")
		return
	}

	fil, escrow, err := s.BankService.Balance(address.String())
	if err != nil {
		s.JSON(w, r, http.StatusInternalServerError, err)
		return
	}

	s.JSON(w, r, http.StatusOK, envelope{"fil": fil, "escrow": escrow})
}

func (s *Server) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	var params AuthorizeParams

	address, ok := r.Context().Value(CtxKeyAddress).(types.Address)
	if !ok {
		s.JSON(w, r, http.StatusBadRequest, "failed to parse header address")
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

	auth, err := s.BankService.Authorize(address.String(), params.Amount)
	if err != nil {
		s.JSON(w, r, http.StatusInternalServerError, err)
		return
	}

	s.JSON(w, r, http.StatusOK, envelope{"fil": auth.Available, "escrow": auth.Escrow, "id": auth.UUID})
}

func (s *Server) handleRefund(w http.ResponseWriter, r *http.Request) {
	address, ok := r.Context().Value(CtxKeyAddress).(types.Address)
	if !ok {
		s.JSON(w, r, http.StatusBadRequest, "failed to parse header address")
		return
	}

	balances, err := s.BankService.Refund(address.String())
	if err != nil {
		s.JSON(w, r, http.StatusInternalServerError, err)
		return
	}

	s.JSON(w, r, http.StatusOK, envelope{"fil": balances.Available, "escrow": balances.Escrow, "expired": balances.Expired})
}

func (s *Server) handleRedeem(w http.ResponseWriter, r *http.Request) {
	var params RedeemParams

	address, ok := r.Context().Value(CtxKeyAddress).(types.Address)
	if !ok {
		s.JSON(w, r, http.StatusBadRequest, "failed to parse header address")
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

	balances, err := s.BankService.Redeem(address.String(), params.UUID, params.Amount)
	if err != nil {
		s.JSON(w, r, http.StatusInternalServerError, err)
		return
	}

	s.JSON(w, r, http.StatusOK, envelope{"excess": balances.Excess, "sp": balances.SP, "cli": balances.CLI})
}

func (s *Server) handleVerify(w http.ResponseWriter, r *http.Request) {
	var params VerifyParams

	address, ok := r.Context().Value(CtxKeyAddress).(types.Address)
	if !ok {
		s.JSON(w, r, http.StatusBadRequest, "failed to parse header address")
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

	err := s.BankService.Verify(address.String(), params.UUID, params.Amount)
	if err != nil {
		s.JSON(w, r, http.StatusInternalServerError, err)
		return
	}

	s.JSON(w, r, http.StatusOK, envelope{"authorization": "valid"})
}
