package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"garago/internal/apis"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := apis.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
