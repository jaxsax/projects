package httpserver

import (
	"context"

	"github.com/go-logr/logr"
)

type Options struct {
	HTTPAddress string `long:"http_address" description:"address to listen for http requests on" default:"0.0.0.0:8080" env:"HTTP_ADDRESS"`
}

type Server struct {
	opts   *Options
	logger logr.Logger
}

func New(opts *Options, logger logr.Logger) *Server {
	return &Server{
		opts:   opts,
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.V(0).Info("starting")
	return nil
}
func (s *Server) Stop(ctx context.Context) error {
	return nil
}
