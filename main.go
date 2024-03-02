package main

import (
	"context"
	"embed"
	"log/slog"
	"net"
	"os"

	"github.com/fabiodcorreia/despensa-app/internal/handlers"
	"github.com/fabiodcorreia/despensa-app/internal/routes"
	"github.com/fabiodcorreia/despensa-app/internal/server"
	"github.com/fabiodcorreia/despensa-app/internal/services"
	"github.com/fabiodcorreia/despensa-app/internal/storage"
)

// default name placeholder
var name = "app"

// default version placeholder
var version = "dev"

//go:embed database/migration/*.sql
var migrations embed.FS

//go:embed public
var public embed.FS

func main() {
	slog.Info("Starting application", "name", name, "version", version)

	if err := run(); err != nil {
		slog.Error("App exited with error", "error", err)
		os.Exit(1)
	}
	slog.Info("App exited with success")
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store, err := storage.NewStoreWithMigrations(
		ctx,
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

	s := server.NewServer(
		ctx,
		server.WithSecure(),
		server.WithRecover(),
		server.WithRemoveTrailingSlash(),
		server.WithLogger(),
		// server.WithSession(),
	)

	s.AddPublic(public)

	s.AddRoutes(routes.HomeRoute{})
	s.AddRoutes(routes.NewLocationRoute(handlers.NewLocationHandler(services.NewLocationService(store))))
	s.AddRoutes(routes.NewItemRoute(handlers.NewItemHandler(services.NewItemService(store))))
	// s.AddRoutes(routes.SearchRoute{})

	go s.Start(net.JoinHostPort(server.GetAddress(), server.GetPort()))
	slog.Info("Server started", "address", server.GetAddress(), "port", server.GetPort())

	return s.WaitAndTerminate()
}
