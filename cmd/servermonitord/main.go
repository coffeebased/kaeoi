package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coffeebased/kaeoi/internal/gameserver"
	"github.com/coffeebased/kaeoi/internal/httpapi"
	"github.com/coffeebased/kaeoi/internal/sqlite"
	"github.com/coffeebased/kaeoi/pkg/poll"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := sqlite.OpenDB("servermonitordb.sqlite")
	if err != nil {
		logger.Error("open sqlite db", "err", err)
		return 1
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("close sqlite db", "err", err)
		}
	}()

	store, err := sqlite.NewStore(db)
	if err != nil {
		logger.Error("server config store", "err", err)
		return 1
	}

	servers, err := store.List(ctx)
	if err != nil {
		logger.Error("store list", "err", err)
		return 1
	}

	targets := gameserver.TargetsFromServers(servers)

	poller, err := poll.New(poll.Options{})
	if err != nil {
		logger.Error("new poller", "err", err)
		return 1
	}

	monitor := gameserver.NewMonitor(targets, poller)

	go monitor.Run(ctx)

	handler := httpapi.New(ctx, logger, store, monitor, 10*time.Second)

	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErr := make(chan error, 1)

	go func() {
		logger.Info("starting http server", "addr", httpServer.Addr)
		serverErr <- httpServer.ListenAndServe()
	}()

	exitCode := 0
	shouldShutdown := true

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")

	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http listen and serve", "err", err)
			exitCode = 1
		}

		shouldShutdown = false
	}

	if shouldShutdown {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("http shutdown", "err", err)
			exitCode = 1
		} else {
			logger.Info("http server stopped")
		}
	}

	return exitCode
}
