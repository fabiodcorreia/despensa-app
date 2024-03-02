package server

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/fabiodcorreia/despensa-app/internal/routes"
	"github.com/labstack/echo/v4"
)

type Server struct {
	serverCtx context.Context
	stopCtx   context.CancelFunc
	engine    *echo.Echo
}

func NewServer(ctx context.Context, options ...ServerOption) *Server {
	engine := echo.New()
	engine.HideBanner = true
	engine.HidePort = true

	serverCtx, stopCtx := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	s := &Server{
		serverCtx,
		stopCtx,
		engine,
	}

	for _, opt := range options {
		opt(s)
	}

	return s
}

func (s *Server) AddPublic(assets fs.FS) {
	s.engine.StaticFS("/", echo.MustSubFS(assets, "public"))
}

func (s *Server) AddRoutes(routes routes.Routes) {
	for _, r := range routes.Routes() {
		s.engine.Add(r.Method, r.Path, r.Handle)
		slog.Info("Added route", "method", r.Method, "path", r.Path)
	}
}

// Start the server and should be called in with go routine
//
// go s.Start("127.0.0.1:8080")
func (s *Server) Start(bindAddr string) {
	if err := s.engine.Start(bindAddr); err != http.ErrServerClosed {
		// s.Log().Errorf("server start : %v", err)
		slog.Error("Server starting", err)
		s.stopCtx() // If error on start call stop and exit
	}
}

// WaitAndTerminate will wait until a termination signal is received and
// then shutdown the server.
func (s *Server) WaitAndTerminate() error {
	s.wait()
	return s.terminate()
}

// Wait will wait until a interrupt signal is received
func (s *Server) wait() {
	defer s.stopCtx()
	// Wait for interrupt signal to gracefully shutdown the server.
	<-s.serverCtx.Done()
}

// Stop will tell the server to shutdown with a 10s timeout
func (s *Server) terminate() error {
	ctx, cancel := context.WithTimeout(s.serverCtx, 30*time.Second)
	defer cancel()
	return s.engine.Shutdown(ctx)
}
