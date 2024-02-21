package main

import (
	"embed"
	"net/http"

	"github.com/fabiodcorreia/despensa-app/internal/handlers"
	"github.com/fabiodcorreia/despensa-app/internal/server"
	"github.com/fabiodcorreia/despensa-app/internal/services"
	"github.com/fabiodcorreia/despensa-app/internal/storage"
)

// default version placeholder
var version = "dev"

//go:embed database/migration/*.sql
var migrations embed.FS

//go:embed public/dist
var public embed.FS

func main() {

	s := server.NewServer()
	s.Log("Despensa server version " + version)

	store, err := storage.NewStoreWithMigrations("despensa.db", migrations)
	if err != nil {
		s.Log(err.Error())
		panic(err)
	}

	if err := store.Connect(); err != nil {
		s.Log(err.Error())
		panic(err)
	}
	defer store.Disconnect()

	s.WithPublic(public)

	searchService := services.NewSearch(store)
	searchHandler := handlers.NewSearch(searchService)
	s.AddRoute(http.MethodGet, "/search", searchHandler.View)
	s.AddRoute(http.MethodPost, "/search", searchHandler.Search)

	locationService := services.NewLocation(store)
	locationHandler := handlers.NewLocation(locationService)

	s.AddRoute(http.MethodGet, "/location/:id", locationHandler.View)

	go s.Start("127.0.0.1:8080")
	s.Log("server ready at http://localhost:8080")
	s.WaitAndTerminate()
	s.Log("server shutdown completed")
}

/*
https://templ.guide/project-structure/project-structure
1. Setup the store
  - Will handle the storage and retrieval of the data
2. Setup the services with the store
  - Will handle the business logic and coordination with the store
3. Setup the handler with the service
  - Will receive the http request, calls the service and renders the component
*/
