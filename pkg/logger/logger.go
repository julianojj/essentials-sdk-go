package logger

import (
	"log/slog"
	"os"
)

type (
	Logger interface {
		Info(message string, data map[string]any)
		Error(message string, err error)
	}
	Slog struct {
		logger *slog.Logger
	}
)

var _ Logger = (*Slog)(nil)

func NewSlog() *Slog {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
	logger := slog.New(logHandler)
	return &Slog{
		logger,
	}
}

func (s *Slog) Info(message string, data map[string]any) {
	s.logger.Info(
		message,
		"data", data,
	)
}

func (s *Slog) Error(message string, err error) {
	s.logger.Error(
		message,
		"err", err,
	)
}
