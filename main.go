package main

import (
	"log/slog"
	"os"

	"github.com/eadydb/nebulae/pkg/utils/log"
)

func main() {
}

func init() {
	opts := log.PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := log.NewPrettyHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
