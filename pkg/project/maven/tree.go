package maven

import (
	"errors"
	"log/slog"
)

type DependencyTree struct{}

// Dependency represents a maven dependency
type Dependency struct {
	GroupId    string // group id
	ArtifactId string // artifact id
	Version    string // version
	Scope      string // scope
}

func parserMavenDependencyTree(fileName string) ([]DependencyTree, error) {
	if fileName == "" {
		return nil, errors.New("filename is empty")
	}
	slog.Info("parsing maven dependency tree ", slog.String("filename", fileName))
	return nil, nil
}
