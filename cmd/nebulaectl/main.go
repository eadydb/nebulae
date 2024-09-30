package main

import (
	"log/slog"
	"os"

	"github.com/eadydb/nebulae/pkg/cmd/nebulaectl"
	"github.com/eadydb/nebulae/pkg/config"
	"github.com/eadydb/nebulae/pkg/utils/log"
	"github.com/eadydb/nebulae/pkg/utils/signals"
	"github.com/spf13/pflag"
)

func main() {
	flagset := pflag.NewFlagSet("global", pflag.ContinueOnError)
	flagset.ParseErrorsWhitelist.UnknownFlags = true
	flagset.Usage = func() {}

	ctx := signals.SetupSignalContext()
	ctx, err := config.InitFlags(ctx, flagset)
	if err != nil {
		_, _ = os.Stderr.Write([]byte(flagset.FlagUsages()))
		slog.Error("Init config flags", slog.String("err", err.Error()))
		os.Exit(1)
	}

	command := nebulaectl.NewCommand(ctx)
	command.PersistentFlags().AddFlagSet(flagset)
	err = command.ExecuteContext(ctx)
	if err != nil {
		slog.Error("Execute exit", slog.String("err", err.Error()))
		os.Exit(1)
	}
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
