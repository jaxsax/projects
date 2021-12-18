package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/felixge/httpsnoop"

	"github.com/jaxsax/projects/tapeworm/botv2/internal"
	"github.com/jaxsax/projects/tapeworm/botv2/links"

	"go.uber.org/zap"
)

type Server struct {
	*zap.Logger
	cfg             *internal.Config
	linksRepository links.Repository
	staticDirPath   string
}

func NewServer(
	logger *zap.Logger,
	cfg *internal.Config,
	linksRepository links.Repository,
	staticDirPath string,
) *Server {
	return &Server{
		Logger:          logger,
		cfg:             cfg,
		linksRepository: linksRepository,
		staticDirPath:   staticDirPath,
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

func (s *Server) Run() error {
	s.Info("listening", zap.Int("port", s.cfg.Port))

	if s.staticDirPath != "" {
		s.Info("launching with static dirs", zap.String("dir", s.staticDirPath))
		fs := http.FileServer(http.Dir(s.staticDirPath))
		http.Handle("/", fs)
	}

	http.Handle("/api/links", Gzip(s.LoggerMiddleware(s.apiLinks())))
	return http.ListenAndServe(fmt.Sprintf(":%v", s.cfg.Port), nil)
}
