package main

import (
	"context"
	"embed"
	"log/slog"
	"net"
	"os"

	"github.com/fabiodcorreia/despensa-app/internal/routes"
	"github.com/fabiodcorreia/despensa-app/internal/server"
	"github.com/fabiodcorreia/despensa-app/internal/storage"
)

var name = "app"

// default version placeholder
var version = "dev"

//go:embed database/migration/*.sql
var migrations embed.FS

//go:embed public
var public embed.FS

func main() {
	ctx := context.Background()
	slog.Info("Starting application", "name", name, "version", version)

	if err := run(ctx); err != nil {
		slog.Error("App exited with error", err)
		os.Exit(1)
	}
	slog.Info("App exited with success")
}

func run(ctx context.Context) error {
	store, err := storage.NewStoreWithMigrations(
		server.GetDatabaseFile(),
		migrations,
	)
	if err != nil {
		return err
	}
	slog.Info("Store initialized and migrations executed")

	if err = store.Connect(); err != nil {
		return err
	}
	defer store.Disconnect()
	slog.Info("Store connected", "database", server.GetDatabaseFile())

	// https://github.com/samber/slog-echo/tree/main
	// s := server.NewServer(ctx)
	s := server.NewServer(
		ctx,
		server.WithSecure(),
		server.WithRecover(),
		server.WithRemoveTrailingSlash(),
		server.WithLogger(),
	)
	s.AddPublic(public)

	homeRoutes := routes.Home{}
	s.AddRoutes(homeRoutes.Routes()...)

	go s.Start(net.JoinHostPort(server.GetAddress(), server.GetPort()))
	slog.Info("Server started", "address", server.GetAddress(), "port", server.GetPort())

	return s.WaitAndTerminate()
}
