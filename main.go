package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	cmd "weathersvc/app/cmd"
)

func main() {
	// signal notifycontext allows for capturing shut down/service stop signals to allow graceful stopping of service
	ctx, stop := signal.NotifyContext(context.Background(), []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}...)
	defer stop()
	if err := cmd.Execute(ctx); err != nil {
		log.Printf("Failed to start server, err: %v", err)
	}
}
