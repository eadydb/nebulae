package maven

import (
	"log/slog"
	"testing"
)

func TestParseMavenDependencyTxt(t *testing.T) {
	deps, err := parseMavenDependencyTxt("deps.txt")
	if err != nil {
		t.Error(err)
	}
	slog.Info("maven dependency tree", slog.Any("deps", deps))
}
