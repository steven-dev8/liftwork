package main

import (
	"context"
	"errors"
	"liftwork/internal/api"
	"liftwork/internal/config"
	"liftwork/internal/database"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httpapi "liftwork/internal/http"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("api stopped", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := database.Open(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()

	application := api.New(pool, &cfg)

	server := &http.Server{
		Addr: cfg.HTTPAddr,
		Handler: httpapi.NewRouter(
			logger,
			&cfg,
			application.Handlers,
		),
		ReadHeaderTimeout: cfg.HTTPReadHeaderTimeout,
	}

	serverError := make(chan error, 1)
	go func() {
		logger.Info("api listening", "address", cfg.HTTPAddr)
		serverError <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return nil
	case err := <-serverError:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
