package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
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

	queryParamDecoder *schema.Decoder
}

func New(opts *Options, s *db.Store, logger logr.Logger) *Server {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	return &Server{
		opts:              opts,
		store:             s,
		logger:            logger,
		queryParamDecoder: decoder,
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
		ctx := logging.WithContext(r.Context())

		if err := r.ParseForm(); err != nil {
			respondWithError(ctx, w, http.StatusBadRequest, "Failed to parseform")
			return
		}

		type request struct {
			Page  int
			Limit int
			Query string
		}

		type response struct {
			Links        []*types.Link `json:"links"`
			Total        int           `json:"total"`
			ItemsPerPage int           `json:"items_per_page"`
			Page         int           `json:"page"`
		}

		var req request
		err := s.queryParamDecoder.Decode(&req, r.Form)
		if err != nil {
			respondWithError(ctx, w, http.StatusBadRequest, "Failed to decode form")
			return
		}

		if req.Limit == 0 {
			req.Limit = 15
		}

		if req.Page == 0 {
			req.Page = 1
		}

		req.Limit = int(math.Min(float64(req.Limit), float64(100)))

		filter := &types.LinkFilter{
			PageNumber: req.Page,
			Limit:      req.Limit,
		}

		if req.Query != "" {
			filter.TitleSearch = req.Query
		}

		links, err := s.store.ListLinksWithFilter(ctx, filter)
		if err != nil {
			logr.FromContextOrDiscard(ctx).Error(err, "failed to retrieve links")
			respondWithError(ctx, w, http.StatusInternalServerError, "Failed to retrieve links")
			return
		}

		countFilter := filter
		countFilter.Limit = 0
		totalCount, err := s.store.CountLinksWithFilter(ctx, countFilter)
		if err != nil {
			logr.FromContextOrDiscard(ctx).Error(err, "failed to count links")
			respondWithError(ctx, w, http.StatusInternalServerError, "Failed to retrieve links")
			return
		}

		resp := &response{
			Links:        links,
			Total:        totalCount,
			ItemsPerPage: req.Limit,
			Page:         req.Page,
		}
		respondWithJSON(ctx, w, http.StatusOK, resp)
	}).Methods(http.MethodGet)

	m.HandleFunc("/api/links/get", func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.WithContext(r.Context())

		linkByUrl := r.URL.Query().Get("url")
		if linkByUrl == "" {
			respondWithError(ctx, w, http.StatusBadRequest, "Invalid url")
			return
		}

		linkWithoutScheme := linkByUrl
		if strings.Contains(linkByUrl, "://") {
			parts := strings.Split(linkByUrl, "://")
			if len(parts) < 2 {
				respondWithError(ctx, w, http.StatusBadRequest, "Invalid url, failed to remove scheme")
				return
			}

			linkWithoutScheme = parts[1]
		}

		links, err := s.store.ListLinksWithFilter(ctx, &types.LinkFilter{
			LinkWithoutScheme: linkWithoutScheme,
		})
		if err != nil {
			respondWithError(ctx, w, http.StatusInternalServerError, "Failed to list links")
			return
		}

		type response struct {
			Links []*types.Link `json:"links"`
		}

		resp := &response{
			Links: links,
		}
		respondWithJSON(ctx, w, http.StatusOK, resp)
	}).Methods(http.MethodGet)

	m.HandleFunc("/api/links/get_by_domain", func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.WithContext(r.Context())

		if err := r.ParseForm(); err != nil {
			respondWithError(ctx, w, http.StatusBadRequest, "Failed to parseform")
			return
		}

		type request struct {
			Domain string
		}

		var req request
		err := s.queryParamDecoder.Decode(&req, r.Form)
		if err != nil {
			respondWithError(ctx, w, http.StatusBadRequest, "Failed to decode form")
			return
		}

		links, err := s.store.ListLinksWithFilter(ctx, &types.LinkFilter{
			Domain:     req.Domain,
			UniqueLink: true,
		})
		if err != nil {
			respondWithError(ctx, w, http.StatusInternalServerError, "Failed to list links")
			return
		}

		type response struct {
			Links []*types.Link `json:"links"`
		}

		resp := &response{
			Links: links,
		}
		respondWithJSON(ctx, w, http.StatusOK, resp)
	})

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
		cacheBust := func(h http.Handler) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/" || r.URL.Path == "index.html" {
					w.Header().Set("Cache-Control", "max-age=0, must-revalidate")
				}
				h.ServeHTTP(w, r)
			}
		}
		m.PathPrefix("/").HandlerFunc(cacheBust(fs))
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
