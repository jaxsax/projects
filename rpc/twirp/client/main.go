package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jaxsax/projects/rpc/twirp/helloworld"
)

func main() {
	c := helloworld.NewHelloWorldJSONClient("http://localhost:9090", &http.Client{})

	r, err := c.Hello(context.TODO(), &helloworld.HelloReq{Subject: "John"})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("response: %v\n", r.Text)
}
