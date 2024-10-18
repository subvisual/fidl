package proxy

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type envelope map[string]any

func (s *Server) Routes(r chi.Router) {
	r.Get("/proxy-routes", s.handleProxy)
}

func (s *Server) handleProxy(w http.ResponseWriter, r *http.Request) {
	s.JSON(w, r, http.StatusOK, envelope{"proxy": "TODO: needed routes and handlers"})
}
