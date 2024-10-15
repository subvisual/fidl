package http

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

func (s *Server) LogError(r *http.Request, err error) {
	rid := middleware.GetReqID(r.Context())
	s.Log.Error(err.Error(), zap.String("rid", rid), zap.String("method", r.Method), zap.String("path", r.URL.Path))
}

func (s *Server) LogWarn(r *http.Request, err error) {
	rid := middleware.GetReqID(r.Context())
	s.Log.Warn(err.Error(), zap.String("rid", rid), zap.String("method", r.Method), zap.String("path", r.URL.Path))
}

func (s *Server) LogDebug(r *http.Request, err error) {
	rid := middleware.GetReqID(r.Context())
	s.Log.Debug(err.Error(), zap.String("rid", rid), zap.String("method", r.Method), zap.String("path", r.URL.Path))
}

func FormatValidationErrors(err error) map[string]string {
	//nolint
	errs := err.(validator.ValidationErrors)
	payload := make(map[string]string)
	for _, fe := range errs {
		payload[strings.ToLower(fe.Field())] = fe.ActualTag()
	}

	return payload
}

func GetHTTPStatusFromStoreError(err error) (int, string) {
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound, "404 Not Found"
	}

	var pqError *pq.Error
	if errors.As(err, &pqError) {
		if strings.Contains(pqError.Error(), "unique constraint") {
			return http.StatusConflict, "409 Conflict"
		} else if strings.Contains(pqError.Error(), "violates check constraint") {
			return http.StatusUnprocessableEntity, "422 Unprocessable Entity"
		}
	}

	return http.StatusInternalServerError, "500 Internal Server Error"
}
