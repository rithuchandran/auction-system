package main

import (
	"bidder/pkg/contract"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"

	"bidder/pkg/config"
	"bidder/pkg/handler"
	"bidder/pkg/service"
)

func main() {
	appConfig, err := config.New()
	if err != nil {
		log.Fatalf("Failed to initialise config: %v", err)
	}

	bidderID := uuid.New().String()
	bidService := &service.Bidder{
		ID:    bidderID,
		Delay: appConfig.Delay,
	}
	router := handler.New(handler.NewBidHandler(bidService))
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", appConfig.AppPort),
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on Port: %v", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	go func() {
		client := &http.Client{}
		b := &bytes.Buffer{}
		registerRequest := contract.RegisterRequest{
			URL:      appConfig.AppURL,
			Port:     appConfig.AppPort,
			BidderID: bidderID,
		}
		if err = json.NewEncoder(b).Encode(registerRequest); err != nil {
			panic(err)
		}
		var req *http.Request
		if req, err = http.NewRequest(http.MethodPost, appConfig.RegistrationURL, b); err != nil {
			panic(err)
		}
		registerBidder(client, req)
	}()
	gracefulStop(server)
}

func registerBidder(client *http.Client, req *http.Request) {
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Error while registering app: %v", err)
		registerBidder(client, req)
	}
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
