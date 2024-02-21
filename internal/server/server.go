package server

import (
	"context"
	"io/fs"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/fabiodcorreia/despensa-app/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Server struct {
	ctx  context.Context
	stop context.CancelFunc
	svr  *echo.Echo
}

func NewServer() *Server {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	svr := echo.New()
	svr.HideBanner = true
	svr.HidePort = true

	svr.Logger.SetLevel(log.INFO)
	svr.Logger.SetPrefix("despensa")
	svr.Use(middleware.Secure())
	svr.Use(middleware.Recover())
	svr.Use(middleware.Logger())
	svr.Pre(middleware.RemoveTrailingSlash())
	// https://echo.labstack.com/docs/middleware/cors

	svr.HTTPErrorHandler = handlers.HTTPErrorHandler
	return &Server{
		ctx,
		stop,
		svr,
	}
}

func (s *Server) Log(message string) {
	s.svr.Logger.Info(message)
}

// Enable the public folder for static content
func (s *Server) WithPublic(public fs.FS) {
	s.svr.StaticFS("/dist", echo.MustSubFS(public, "public/dist"))
}

func (s *Server) AddRoute(method string, path string, handler echo.HandlerFunc) {
	s.svr.Add(method, path, handler)
}

// Run with go routine
func (s *Server) Start(bindAddr string) {
	if err := s.svr.Start(bindAddr); err != http.ErrServerClosed {
		s.svr.Logger.Error("server start fail with error: " + err.Error())
		s.stop() // If error on start call stop and exit
	}
}

// Wait until a termination signal is received and then shutdown the server.
//
// The shutdown will have a timeout of 10s.
func (s *Server) WaitAndTerminate() error {
	s.wait()
	return s.terminate()
}

// Wait will wait until a interrupt signal is received
func (s *Server) wait() {
	defer s.stop()
	// Wait for interrupt signal to gracefully shutdown the server.
	<-s.ctx.Done()
}

// Stop will tell the server to shutdown with a 10s timeout
func (s *Server) terminate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.svr.Shutdown(ctx)
}
