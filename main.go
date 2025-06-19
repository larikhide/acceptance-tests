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
		Addr:    ":8080",
		Handler: http.HandlerFunc(acceptancetests.SlowHandler),
	}
	server := gracefulshutdown.NewServer(httpServer)

	err := server.ListenAndServe(ctx)
	if err != nil {
		log.Fatalf("fail shutdown gracefully: %v", err)
	}

	log.Println("shutdown gracefully!")
}
