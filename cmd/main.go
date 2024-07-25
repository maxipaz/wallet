package main

import (
	"context"
	"github.com/maxipaz/wallet/cmd/command"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutdownTimeout = 5

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	defer func() {
		slog.Debug("closing application")
		signal.Stop(shutdownSignal)
		cancel()
	}()

	go func() {
		<-shutdownSignal
		cancel()
		time.Sleep(shutdownTimeout * time.Second)
		slog.Debug("fallback exit")
		os.Exit(1)
	}()

	if err := command.NewRootCommand(ctx).Execute(); err != nil {
		os.Exit(1)
	}
}
