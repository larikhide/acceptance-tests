package main

import (
	"context"
	"log"
	"net/http"

	gracefulshutdown "github.com/quii/go-graceful-shutdown"
	"github.com/quii/go-graceful-shutdown/acceptancetests"
)

func main() {
	ctx := context.Background()
	httpServer := &http.Server{
		Addr:    "8082",
		Handler: http.HandlerFunc(acceptancetests.SlowHandler),
	}
	server := gracefulshutdown.NewServer(httpServer)

	if err := server.ListenAndServe(ctx); err != nil {
		log.Fatalf("cannot finish gracefully: %v", err)
	}

	log.Println("shutdowned gracefully")
}
