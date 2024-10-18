package bank

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type envelope map[string]any

func (s *Server) Routes(r chi.Router) {
	r.Post("/register", s.handleRegisterProxy)
	r.Post("/deposit", s.handleDeposit)
	r.Post("/withdraw", s.handleWithdraw)
	r.Get("/balance", s.handleGetBalance)
	r.Post("/authorize", s.handleAuthorize)
	r.Post("/redeem", s.handleRedeem)
}

func (s *Server) handleRegisterProxy(w http.ResponseWriter, r *http.Request) {
	s.JSON(w, r, http.StatusOK, envelope{"bank": "TODO: register proxy"})
}

func (s *Server) handleDeposit(w http.ResponseWriter, r *http.Request) {
	s.JSON(w, r, http.StatusOK, envelope{"bank": "TODO: deposit"})
}

func (s *Server) handleWithdraw(w http.ResponseWriter, r *http.Request) {
	s.JSON(w, r, http.StatusOK, envelope{"bank": "TODO: withdraw"})
}

func (s *Server) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	s.JSON(w, r, http.StatusOK, envelope{"bank": "TODO: get balance"})
}

func (s *Server) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	s.JSON(w, r, http.StatusOK, envelope{"bank": "TODO: authorize"})
}

func (s *Server) handleRedeem(w http.ResponseWriter, r *http.Request) {
	s.JSON(w, r, http.StatusOK, envelope{"bank": "TODO: redeem"})
}
