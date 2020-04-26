package main

import (
	"auctioneer/pkg/handler"
	"auctioneer/pkg/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const defaultAppPort = 8080

func main() {
	service := &service.Auctioneer{
		Bidders: make(map[string]string),
		Client:  &http.Client{},
	}
	auctionHandler := handler.NewAuctionHandler(service)
	router := handler.New(auctionHandler)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", defaultAppPort),
		Handler: router,
	}
	go func() {
		log.Printf("Starting server on Port: %v", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	gracefulStop(server)
}

//listens for quit, terminate and interrupt signals and shuts the server gracefully without interrupting any active connections
func gracefulStop(server *http.Server) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-stop

	log.Printf("Shutting the server down...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		log.Printf("Server stopped")
	}
}
