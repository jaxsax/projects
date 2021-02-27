package main

import (
	"context"
	"net/http"

	"github.com/jaxsax/projects/rpc/twirp/helloworld"
)

type Server struct{}

func (s *Server) Hello(context.Context, *helloworld.HelloReq) (*helloworld.HelloResp, error) {
	return &helloworld.HelloResp{
		Text: "hello",
	}, nil
}

func main() {
	server := &Server{}
	twirpHandler := helloworld.NewHelloWorldServer(server)

	err := http.ListenAndServe(":9090", twirpHandler)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
