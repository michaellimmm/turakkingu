package main

import (
	"context"
	"github/michaellimmm/turakkingu/internal/adapter"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/repository"
	"github/michaellimmm/turakkingu/internal/usecase"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := core.NewConfig()
	if err != nil {
		slog.Error("failed to get config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	repo, err := repository.NewRepo(config)
	if err != nil {
		slog.Error("failed to initialize repository", slog.String("error", err.Error()))
		os.Exit(1)
	}

	uc := usecase.NewUseCase(config, repo)
	server := adapter.NewAdapter(config, uc)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	serverErrChan := make(chan error, 1)

	go func() {
		err = server.Run()
		if err != nil {
			slog.Error("error", slog.String("error", err.Error()))
			serverErrChan <- err
		}
	}()

	slog.Info("server is running....")

	select {
	case sig := <-sigChan:
		slog.Info("shutdown signal received, starting graceful shutdown", slog.String("signal", sig.String()))
		gracefulShutdown(server, repo)
	case err := <-serverErrChan:
		slog.Error("server failed to start", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("shutting down services")
}

func gracefulShutdown(server adapter.AdapterCloser, repo repository.RepoCloser) {
	slog.Info("stopping server...")
	if err := server.Close(context.Background()); err != nil {
		slog.Error("server shutdown error", slog.String("error", err.Error()))
	}

	slog.Info("closing database connections...")
	if err := repo.Close(context.Background()); err != nil {
		slog.Error("repository close error", slog.String("error", err.Error()))
	}
}
