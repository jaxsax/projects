package web

import (
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

func (s *Server) Run() error {
	s.Message(fmt.Sprintf("listening on %v", s.cfg.Port))
	return http.ListenAndServe(fmt.Sprintf(":%v", s.cfg.Port), nil)
}
