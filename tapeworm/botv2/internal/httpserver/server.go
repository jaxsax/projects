package httpserver

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
)

type Options struct {
	HTTPAddress string `long:"http_address" description:"address to listen for http requests on" default:"0.0.0.0:8081" env:"HTTP_ADDRESS"`
}

type Server struct {
	opts   *Options
	logger logr.Logger

	httpServer *http.Server
	store      *db.Store
}

func New(opts *Options, s *db.Store, logger logr.Logger) *Server {
	return &Server{
		opts:   opts,
		store:  s,
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.V(0).Info("starting", "addr", s.opts.HTTPAddress)

	s.httpServer = &http.Server{
		Addr:    s.opts.HTTPAddress,
		Handler: s.buildMux(),
	}

	if err := s.httpServer.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}

	return nil
}

func (s *Server) buildMux() *mux.Router {
	m := mux.NewRouter().StrictSlash(true)

	m.HandleFunc("/liveness", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	return m
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.V(0).Info("shutting down")
	return s.httpServer.Shutdown(ctx)
}
