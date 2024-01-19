package main

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	"goservice/internal/env"
	"goservice/internal/server"
	"log"
)

var ctx = context.Background()

func main() {
	config := &env.Config{}

	if err := envconfig.Process(ctx, config); err != nil {
		log.Fatalln("Fatal Error: Parsing OS ENV")
	}
	if err := server.Serve(config); err != nil {
		log.Fatalln("Fatal Error: Start Up gRPC server")
	}
}
