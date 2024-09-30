package gitlab

import (
	"context"
	"log/slog"
	"testing"

	"github.com/eadydb/nebulae/pkg/repository"
)

func TestGitlabRepository(t *testing.T) {
	ctx := context.Background()
	initDriver(ctx)
	gitlab := NewGitlab(ctx, "gitlab_privateToken", "gitlab_url")
	projects := &GitlabProjects{
		Simple:     true,
		Membership: false,
		Page:       1,
		PerPage:    10,
		Recursive:  false,
	}
	err := gitlab.ScanGitlabHub(projects)
	if err != nil {
		t.Error(err)
	}
}

func TestCloneRepository(t *testing.T) {
	ctx := context.Background()
	initDriver(ctx)
	repository, err := newRepositoryService(ctx).GetRepositoryById(4324)
	if err != nil {
		slog.Error("get repository by id failed", slog.String("err", err.Error()))
		return
	}
	gitlab := NewGitlab(ctx, "gitlab_privateToken", "gitlab_url")
	err = gitlab.Clone(repository, "~/Workspace/Gitlab/eadydb")
	if err != nil {
		slog.Error("clone repository failed", slog.String("err", err.Error()))
	}
}

func initDriver(ctx context.Context) {
	db, err := repository.NewMysqlDriver("", "127.0.0.1", "admin", "HdaTaV2Vf", "payment_test", 18248).GetGromConn(ctx)
	if err != nil {
		panic(err)
	}
	repository.GormConn = db
}
