package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/term"
)

// InitLogger initializes the logger. If the log-format is not specified, it will default to JSON if the output is not a terminal.
// Required viper variables: log-format as string and debug as bool
func InitLogger() {
	cfg := Config{
		Debug:  viper.GetBool("debug"),
		Format: viper.GetString("log-format"),
	}
	InitLoggerWithConfig(cfg)
}

type Config struct {
	Debug  bool
	Format string
}

// InitLoggerWithConfig initializes the logger with the specified configuration.
func InitLoggerWithConfig(cfg Config) {
	isTerminal := term.IsTerminal(int(os.Stdout.Fd()))

	var opts *slog.HandlerOptions
	if cfg.Debug {
		opts = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	}

	switch cfg.Format {
	case "json":
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opts)))
	case "logfmt":
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))
	case "":
		// Default to JSON if not a terminal
		if !isTerminal {
			slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opts)))
		} else {
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))
		}
	default:
		slog.Error("Invalid log format specified")
		os.Exit(1)
	}
}
