package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jaxsax/projects/tapeworm/botv2/internal"
	"github.com/jaxsax/projects/tapeworm/botv2/links"
)

type Server struct {
	*internal.Logger
	cfg             *internal.Config
	linksRepository links.Repository
}

func NewServer(
	logger *internal.Logger,
	cfg *internal.Config,
	linksRepository links.Repository,
) *Server {
	return &Server{
		Logger:          logger,
		cfg:             cfg,
		linksRepository: linksRepository,
	}
}

func (s *Server) handleLinks(w http.ResponseWriter, r *http.Request) {
	dbLinks := s.linksRepository.List()

	resp := struct {
		Links []links.Link
	}{
		Links: dbLinks,
	}

	js, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		s.Log(
			"endpoint", "/api/links",
			"err", err.Error(),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

func (s *Server) Run() error {
	s.Message(fmt.Sprintf("listening on %v", s.cfg.Port))

	http.HandleFunc("/api/links", s.handleLinks)
	return http.ListenAndServe(fmt.Sprintf(":%v", s.cfg.Port), nil)
}
