package gitlab

import (
	"context"
	"testing"
)

func TestGitlabRepository(t *testing.T) {
	gitlab := NewGitlab(context.Background(), "zU1HF3Q9MbYG5J", "https://git.leyaoyao.com")
	projects := &GitlabProjects{
		Simple:     true,
		Membership: false,
		Page:       1,
		PerPage:    10,
	}
	err := gitlab.ScanGitlabHub(projects)
	if err != nil {
		t.Error(err)
	}
}
