package main

import (
	"context"
	"os"
	"os/signal"
	"slurm/go-on-practice-2/http_06/internals/app"
	"slurm/go-on-practice-2/http_06/internals/cfg"

	"github.com/sirupsen/logrus"
)

func main() {
	config := cfg.LoadAndStoreConfig()

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	server := app.NewServer(config, ctx)

	go func() {
		oscall := <-c
		logrus.Printf("system call: %+v", oscall)
		server.Shutdown()
		cancel()
	}()

	server.Serve()

}
