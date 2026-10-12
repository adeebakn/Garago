package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"garago/internal/apis"
)

// main preserves the convenient `go run .` development command.
// The production-oriented executable remains available at ./cmd/api.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := apis.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
