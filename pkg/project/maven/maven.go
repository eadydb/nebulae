package maven

import "context"

type Maven struct {
	Path    string // maven pom.xml path
	Project string // project name
	Ctx     context.Context
}

func NewMaven(ctx context.Context, path, project string) *Maven {
	return &Maven{
		Path:    path,
		Project: project,
		Ctx:     ctx,
	}
}

// Build
func (m *Maven) Build() error {
	return nil
}

func (m *Maven) Dependency() error {
	return nil
}
