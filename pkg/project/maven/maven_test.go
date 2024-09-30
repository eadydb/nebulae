package maven

import (
	"context"
	"log/slog"
	"testing"
)

func TestMavenCli(t *testing.T) {
	ctx := context.Background()
	mvn := NewMaven(ctx, "/Users/eadydb/Workspace/Gitlab/eadydb/charging/back-end/charging-base-starter", "charging-base-starter")
	mvn.FileName = "mvn_deps.txt"
	err := mvn.ExectueCmd(true)
	if err != nil {
		slog.Error("execute maven dependency:tree cmd failed", slog.String("err", err.Error()))
		return
	}
	files, err := mvn.LoadingMvnTxtFile()
	if err != nil {
		slog.Error("loading mvn txt file failed", slog.String("err", err.Error()))
		return
	}
	slog.Info("loading mvn txt file", slog.Any("files", files))
}
