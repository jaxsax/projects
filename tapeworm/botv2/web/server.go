package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/google/uuid"

	"github.com/jaxsax/projects/tapeworm/botv2/internal"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/utils"
	"github.com/jaxsax/projects/tapeworm/botv2/links"

	"go.uber.org/zap"
)

type Server struct {
	*zap.Logger
	cfg             *internal.Config
	linksRepository links.Repository
	staticDirPath   string
	db              *sql.DB
}

func NewServer(
	logger *zap.Logger,
	cfg *internal.Config,
	linksRepository links.Repository,
	staticDirPath string,
	db *sql.DB,
) *Server {
	return &Server{
		Logger:          logger,
		cfg:             cfg,
		linksRepository: linksRepository,
		staticDirPath:   staticDirPath,
		db:              db,
	}
}

func useRequestIDMiddleware(next http.Handler) http.Handler {
	const headerName = "x-request-id"
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		println("@@@@@@")
		id := r.Header.Get(headerName)
		if id == "" {
			id = uuid.New().String()
		}

		ctx := utils.SetRequestID(r.Context(), id)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}

func (s *Server) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := utils.MustGetRequestID(r.Context())
		s.Info("request started",
			zap.String("request.method", r.Method),
			zap.String("request.path", r.URL.String()),
			zap.String("request.id", requestId),
		)

		m := httpsnoop.CaptureMetrics(next, w, r)

		s.Info(
			"request complete",
			zap.Int("response.status_code", m.Code),
			zap.Duration("duration", m.Duration),
			zap.Int64("bytes_written", m.Written),
			zap.String("request.method", r.Method),
			zap.String("request.path", r.URL.String()),
			zap.String("request.id", requestId),
		)
	})
}

func (s *Server) txMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		tx, err := s.db.Begin()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := utils.SetTransaction(r.Context(), tx)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)

		if err := tx.Commit(); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			s.Error("commit failed", zap.Error(err))
			return
		}
	})
}

func (s *Server) apiLinks() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		linkValues, err := links.GetLinks(r.Context())
		if err != nil {
			s.Error("error listing", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := struct {
			Links []*links.Link `json:"links"`
		}{
			Links: linkValues,
		}

		js, err := json.Marshal(resp)
		if err != nil {
			s.Error("error marshalling", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(js)
	}
	return http.HandlerFunc(fn)
}

func middlewareWrapper(handlers ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		for i := len(handlers) - 1; i >= 0; i-- {
			h = handlers[i](h)
		}

		return h
	}
}

func (s *Server) Run() error {
	s.Info("listening", zap.Int("port", s.cfg.Port))

	if s.staticDirPath != "" {
		s.Info("launching with static dirs", zap.String("dir", s.staticDirPath))
		fs := http.FileServer(http.Dir(s.staticDirPath))
		http.Handle("/", fs)
	}

	mwStack := []func(http.Handler) http.Handler{
		useRequestIDMiddleware,
		s.LoggerMiddleware,
		s.txMiddleware,
		Gzip,
	}

	stack := middlewareWrapper(mwStack...)

	http.Handle("/api/links", stack(s.apiLinks()))
	return http.ListenAndServe(fmt.Sprintf(":%v", s.cfg.Port), nil)
}
