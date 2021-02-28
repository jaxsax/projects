package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/felixge/httpsnoop"

	"github.com/jaxsax/projects/tapeworm/botv2/internal"
	"github.com/jaxsax/projects/tapeworm/botv2/links"
	"github.com/jaxsax/projects/tapeworm/botv2/search"

	"go.uber.org/zap"
)

type Server struct {
	*zap.Logger
	cfg             *internal.Config
	linksRepository links.Repository
	searcher        search.LinkSearcher
}

func NewServer(
	logger *zap.Logger,
	cfg *internal.Config,
	linksRepository links.Repository,
	searcher search.LinkSearcher,
) *Server {
	return &Server{
		Logger:          logger,
		cfg:             cfg,
		linksRepository: linksRepository,
		searcher:        searcher,
	}
}

func (s *Server) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(next, w, r)
		internal.Emit(
			s.Logger,
			zap.String("request.method", r.Method),
			zap.String("request.path", r.URL.String()),
			zap.Int("response.status_code", m.Code),
			zap.Duration("duration", m.Duration),
			zap.Int64("bytes_written", m.Written),
		)
	})
}

func (s *Server) apiLinks() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		dbLinks, err := s.linksRepository.List()
		if err != nil {
			s.Logger.Error("error listing", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := struct {
			Links []links.Link
		}{
			Links: dbLinks,
		}

		js, err := json.Marshal(resp)
		if err != nil {
			s.Logger.Error("error marshalling", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(js)
	}
	return http.HandlerFunc(fn)
}

func (s *Server) search() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		q, ok := queryParams["q"]
		if !ok {
			http.Error(w, "no query", http.StatusBadRequest)
			return
		}

		if len(q) == 0 {
			http.Error(w, "no query", http.StatusBadRequest)
			return
		}

		var (
			searchTerm = q[0]
			limit      = 100
			offset     = 0
		)

		s.Logger.Info(
			"sending search request",
			zap.String("term", searchTerm),
			zap.Int("limit", limit),
			zap.Int("offset", offset),
		)
		foundLinks, err := s.searcher.Search(searchTerm, limit, offset)
		if err != nil {
			s.Logger.Error("error searching", zap.Error(err), zap.String("term", searchTerm))
			http.Error(w, "error making search query", http.StatusInternalServerError)
			return
		}

		resp := struct {
			Links []links.Link
		}{
			Links: foundLinks,
		}

		js, err := json.Marshal(resp)
		if err != nil {
			s.Logger.Error("error marshalling", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(js)
	})
}

func (s *Server) Run() error {
	s.Logger.Info("listening")

	http.Handle("/api/links", Gzip(s.LoggerMiddleware(s.apiLinks())))
	http.Handle("/api/search", Gzip(s.LoggerMiddleware(s.search())))
	return http.ListenAndServe(fmt.Sprintf(":%v", s.cfg.Port), nil)
}
