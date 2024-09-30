package gitlab

import (
	"context"
	"log/slog"
	"testing"
)

func TestGetProjectLanguage(t *testing.T) {
	ctx := context.Background()
	initDriver(ctx)
	gitlab := NewGitlab(ctx, "", "")
	language, err := gitlab.GetProjectLanguage("~/Workspace/GitHub/eadydb/nebulae")
	if err != nil {
		slog.Error("get project language failed", slog.String("err", err.Error()))
	}
	slog.Debug("get project language", slog.String("language", language))
}
