package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	mw "github.com/subvisual/fidl/http/middleware"
	"go.uber.org/zap"
)

type Config struct {
	Addr            string
	Fqdn            string
	Port            int
	ListenPort      int
	ReadTimeout     int
	WriteTimeout    int
	ShutdownTimeout int

	Env string
	TLS bool
}

type Server struct {
	listener net.Listener
	server   *http.Server
	router   *chi.Mux
	cfg      *Config
	decoder  *schema.Decoder

	Validate *validator.Validate
	Log      *zap.Logger
}

func New(cfg *Config) *Server {
	srv := &Server{
		server: &http.Server{
			ReadTimeout:       time.Duration(cfg.ReadTimeout) * time.Second,
			WriteTimeout:      time.Duration(cfg.WriteTimeout) * time.Second,
			ReadHeaderTimeout: time.Duration(cfg.ReadTimeout) * time.Second,
		},
		router:   chi.NewRouter(),
		cfg:      cfg,
		decoder:  schema.NewDecoder(),
		Validate: validator.New(),
	}

	srv.server.Handler = srv.router

	return srv
}

func (s *Server) Decode(dst any, src map[string][]string) error {
	// nolint:wrapcheck
	return s.decoder.Decode(dst, src)
}

func (s *Server) Host() string {
	if s.cfg.Port == 80 || s.cfg.Port == 443 {
		return s.cfg.Fqdn
	}

	return fmt.Sprintf("%s:%d", s.cfg.Fqdn, s.cfg.Port)
}

func (s *Server) Schema() string {
	if s.cfg.TLS {
		return "https"
	}

	return "http"
}

func (s *Server) URI() string {
	return fmt.Sprintf("%s://%s", s.Schema(), s.Host())
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.ShutdownTimeout))
	defer cancel()

	//nolint
	return s.server.Shutdown(ctx)
}

func (s *Server) RegisterMiddleWare(middlewares ...func(http.Handler) http.Handler) {
	// Register middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	s.router.Use(mw.NewZap(s.Log))
	s.router.Use(middleware.Recoverer)
	s.router.Use(middlewares...)
}

func (s *Server) RegisterRoutes(routes ...func(chi.Router)) {
	s.router.Route("/api/v1", func(r chi.Router) {
		s.registerHealthCheckRoutes(r)
		for _, routeFn := range routes {
			routeFn(r)
		}
	})
}

func (s *Server) Run() error {
	var err error

	walkFunc := func(
		method string,
		route string,
		_ http.Handler,
		_ ...func(http.Handler) http.Handler) error {
		route = strings.ReplaceAll(route, "/*/", "/")
		s.Log.Info("route", zap.String("method", method), zap.String("path", route))

		return nil
	}

	if err := chi.Walk(s.router, walkFunc); err != nil {
		s.Log.Error("failed to walk routes", zap.Error(err))
	}

	address := fmt.Sprintf("%s:%d", s.cfg.Addr, s.cfg.ListenPort)
	s.listener, err = net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on address: '%s': %w", address, err)
	}

	go func() {
		err = s.server.Serve(s.listener)
	}()

	// nolint
	return err
}
