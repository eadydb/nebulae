package config

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path"

	"github.com/eadydb/nebulae/pkg/consts"
	"github.com/eadydb/nebulae/pkg/utils/envs"
	"github.com/spf13/pflag"
)

type configCtx int

// InitFlags initializes the flags for the configuration.
func InitFlags(ctx context.Context, flags *pflag.FlagSet) (context.Context, error) {
	defaultConfigPath := path.Join(WorkDir, consts.ConfigName)
	config := flags.StringArrayP("config", "c", []string{defaultConfigPath}, "config path")
	_ = flags.Parse(os.Args[1:])

	objs, err := Load(ctx, *config...)
	if err != nil {
		if len(*config) == 1 && (*config)[0] == defaultConfigPath && errors.Is(err, os.ErrNotExist) {
			slog.Debug("Load config", slog.String("path", (*config)[0]), slog.String("err", err.Error()))

			return ctx, nil
		}
		return nil, err
	}

	if len(objs) == 0 {
		slog.Debug("Load config", slog.String("path", (*config)[0]), slog.String("err", "empty config"))

		return ctx, nil
	}

	slog.Debug("Load config", slog.String("path", (*config)[0]))
	return context.WithValue(ctx, configCtx(0), objs), nil
}

var WorkDir = envs.GetEnvWithPrefix("WORKDIR", func() string {
	dir, err := os.UserHomeDir()
	if err != nil || dir == "" {
		return path.Join(os.TempDir(), consts.ProjectName)
	}
	return path.Join(dir, "."+consts.ProjectName)
}())
