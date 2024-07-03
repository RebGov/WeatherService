package app

import (
	"context"
	"log"
	"weathersvc/app/config"
	"weathersvc/app/server"
	"weathersvc/app/service"
)

func Execute(ctx context.Context) error {
	conf, err := config.NewApp(ctx)
	if err != nil {
		return err
	}
	svc, err := service.NewService(ctx, &conf)
	if err != nil {
		return err
	}
	svr := server.NewServer(conf, svc)
	log.Println("Service Starting")
	// listen for context cancellation to handle signal inter
	go func() {
		<-ctx.Done()
		if err := svr.Close(); err != nil {
			log.Printf("Failed to gracefully shutdown the service: %v", err)
		}
	}()
	// Start the HTTP server.
	if err := svr.Open(); err != nil {
		log.Printf("Failed to create rest server: %v", err)
		return err
	}
	return nil
}
