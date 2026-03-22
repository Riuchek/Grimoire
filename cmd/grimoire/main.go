package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"grimoire/internal/app"
	"grimoire/internal/config"
)

func main() {
	cfg := config.Load()
	if cfg.Token == "" {
		log.Fatal("DISCORD_TOKEN is not set")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
