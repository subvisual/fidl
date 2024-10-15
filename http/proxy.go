package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) registerProxyRoutes(r chi.Router) {
	r.Get("/proxy-routes", s.handleProxy)
}

func (s *Server) handleProxy(w http.ResponseWriter, r *http.Request) {
	s.JSON(w, r, http.StatusOK, envelope{"proxy": "TODO: needed routes and handlers"})
}
