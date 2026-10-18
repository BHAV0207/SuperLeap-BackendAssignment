package observability

import (
	"log/slog"
	"os"
	"strings"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/version"
)

// NewLogger sets up a JSON logger on stdout so the logs can be read by
// kubectl logs / Loki without any parsing rules.
func NewLogger(level, service, env string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})

	logger := slog.New(handler).With(
		"service", service,
		"env", env,
		"version", version.Version,
		"commit", version.Commit,
	)

	slog.SetDefault(logger)
	return logger
}
