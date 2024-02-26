package server

import (
	"log/slog"

	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

type ServerOption func(*Server)

func WithSecure() ServerOption {
	return func(s *Server) {
		s.engine.Use(middleware.Secure())
	}
}

func WithRecover() ServerOption {
	return func(s *Server) {
		s.engine.Use(middleware.Recover())
	}
}

func WithRemoveTrailingSlash() ServerOption {
	return func(s *Server) {
		s.engine.Use(middleware.RemoveTrailingSlash())
	}
}

func WithLogger() ServerOption {

	config := slogecho.Config{
		DefaultLevel:     slog.LevelDebug,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithUserAgent:      false,
		WithRequestID:      false,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,
	}

	return func(s *Server) {
		s.engine.Use(slogecho.NewWithConfig(slog.Default(), config))
	}
}
