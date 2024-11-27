package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"github.com/subvisual/fidl/http/jsend"
	"github.com/subvisual/fidl/validation"
)

type envelope map[string]any

func (s *Server) JSON(w http.ResponseWriter, r *http.Request, code int, value any) {
	var status int
	var body any

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//nolint
	if err, ok := value.(error); ok {
		var validationError validator.ValidationErrors
		var pqError *pq.Error
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &validationError):
			status, body = code, FormatValidationErrors(validationError)
		case errors.Is(err, sql.ErrNoRows):
			status, body = GetHTTPStatusFromStoreError(err)
		case errors.As(err, &pqError):
			status, body = GetHTTPStatusFromStoreError(err)
		case errors.As(err, &maxBytesError):
			status, body = http.StatusRequestEntityTooLarge, envelope{"file": "too large", "max": maxBytesError.Limit}
		case errors.Is(err, validation.ErrInvalidContentLength):
			status, body = http.StatusUnprocessableEntity, envelope{"file": "size"}
		case errors.Is(err, validation.ErrInvalidMimeType):
			status, body = http.StatusUnprocessableEntity, envelope{"file": "mime"}
		default:
			status, body = code, err.Error()
		}

		if status < 500 {
			s.LogDebug(r, err)
		} else {
			s.LogError(r, err)
		}
	} else {
		if e, ok := value.(envelope); ok {
			if v, ok := e["err"]; ok {
				//nolint
				s.LogDebug(r, v.(error))
				delete(e, "err")
			}
		}
		status, body = code, value
	}

	var payload jsend.Payload
	switch {
	case status >= 200 && status < 300:
		payload = jsend.Ok(body)
	case status >= 400 && status < 500:
		payload = jsend.Fail(body)
	case status > 500:
		payload = jsend.Fail(body)
	default:
		payload = jsend.Error("The server encountered a problem and could not process your request")
	}

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		s.LogError(r, err)
	}
}

func (s *Server) DecodeJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(destination)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}

			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// Setting DisallowUnknownFields may trigger this error when an unknown field is found.
		// Since there is no specific error type for this case a string match is used.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	}

	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
