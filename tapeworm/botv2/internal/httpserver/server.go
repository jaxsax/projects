package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/types"
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

	m.HandleFunc("/api/links", func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Links []*types.Link `json:"links"`
		}

		links, err := s.store.ListLinks(r.Context())
		if err != nil {
			logging.Logger.Error(err, "failed to list links")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("error"))
			return
		}

		resp := &response{
			Links: links,
		}
		responseBytes, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(responseBytes)
	}).Methods(http.MethodGet)

	return m
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.V(0).Info("shutting down")
	return s.httpServer.Shutdown(ctx)
}
