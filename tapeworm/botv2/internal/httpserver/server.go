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
	HTTPAddress       string `long:"http_address" description:"address to listen for http requests on" default:"0.0.0.0:8081" env:"HTTP_ADDRESS"`
	FrontendAssetPath string `long:"frontend_asset_path" description:"path to frontend build" env:"FRONTEND_ASSET_PATH"`
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
			respondWithError(r.Context(), w, http.StatusInternalServerError, "Failed to retrieve links")
			return
		}

		resp := &response{
			Links: links,
		}
		respondWithJSON(r.Context(), w, http.StatusOK, resp)
	}).Methods(http.MethodGet)

	m.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.WithContext(r.Context())

		query := r.URL.Query().Get("q")
		if query == "" {
			respondWithError(ctx, w, http.StatusBadRequest, "Invalid query")
			return
		}

		resp, err := s.store.Search(ctx, &types.SearchRequest{
			FullText: query,
		})
		if err != nil {
			respondWithError(ctx, w, http.StatusInternalServerError, "Failed to search links")
			return
		}

		respondWithJSON(ctx, w, http.StatusOK, resp)
	}).Methods(http.MethodGet)

	if s.opts.FrontendAssetPath != "" {
		fs := http.FileServer(http.Dir(s.opts.FrontendAssetPath))
		m.PathPrefix("/").Handler(fs)
	}

	return m
}

func respondWithError(ctx context.Context, w http.ResponseWriter, code int, message string) {
	respondWithJSON(ctx, w, code, map[string]string{
		"error": message,
	})
}

func respondWithJSON(ctx context.Context, w http.ResponseWriter, code int, payload interface{}) {
	responseBytes, err := json.Marshal(payload)
	if err != nil {
		logr.FromContextOrDiscard(ctx).
			Error(
				err, "failed to marshal response payload",
				"code", code,
				"payload", payload,
			)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("hi"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(responseBytes)
	if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "failed to write to response to client")
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.V(0).Info("shutting down")
	return s.httpServer.Shutdown(ctx)
}
