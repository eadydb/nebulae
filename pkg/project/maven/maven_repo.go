package maven

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/eadydb/nebulae/pkg/gitlab"
	"github.com/eadydb/nebulae/pkg/project/pom"
	"github.com/eadydb/nebulae/pkg/repository"
	"gorm.io/gorm"
)

type GitlabRepositoryMaven struct {
	Ctx      context.Context
	DB       *gorm.DB
	PageSize int
	Language string
}

func NewGitlabRepositoryMaven(ctx context.Context, pageSize int, language string) *GitlabRepositoryMaven {
	return &GitlabRepositoryMaven{
		PageSize: pageSize,
		Ctx:      ctx,
		DB:       repository.GormConn,
		Language: language,
	}
}

func (g *GitlabRepositoryMaven) Do() error {
	return nil
}

func (g *GitlabRepositoryMaven) DoMvnCommand(repo *gitlab.RepositoryService, limit int) error {
	repos, err := repo.FindRepositories(g.PageSize, limit, g.Language)
	if err != nil {
		slog.Error("find gitlab repository failed", slog.String("err", err.Error()), slog.Int("limit", limit), slog.String("language", g.Language))
		return err
	}
	if len(repos) == 0 {
		return nil
	}
	var wg sync.WaitGroup
	maxConcurrent := 2
	sem := make(chan struct{}, maxConcurrent)
	errChan := make(chan error, len(repos))

	for _, repo := range repos {
		sem <- struct{}{}
		wg.Add(1)
		go func(repo gitlab.GitlabRepository) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("panic in goroutine: %v", r)
				}
				<-sem
			}()
			err := g.DoMvnDependency(repo)
			if err != nil {
				errChan <- err
			}
		}(repo)
	}
	wg.Wait()
	close(errChan)
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		slog.Error("encountered %d errors during execution:", slog.Any("count", len(errs)), slog.Any("errors", errs))
	}
	if len(repos) < g.PageSize {
		return g.DoMvnCommand(repo, limit+1)
	}
	return nil
}

func (g *GitlabRepositoryMaven) DoMvnDependency(repo gitlab.GitlabRepository) error {
	maven := NewMaven(g.Ctx, repo.WorkspaceDir, repo.Name)
	if err := maven.ExectueCmd(true); err != nil {
		slog.Error("execute maven dependency:tree failed", slog.String("project", repo.Name), slog.String("path", repo.Path), slog.String("err", err.Error()))
		return err
	}
	txtPaths, err := maven.LoadingMvnTxtFile()
	if err != nil {
		slog.Error("loading maven txt file failed", slog.String("project", repo.Name), slog.String("path", repo.Path), slog.String("err", err.Error()))
		return err
	}
	for _, txtPath := range txtPaths {
		dependencyTxt, err := parseMavenDependencyTxt(txtPath)
		if err != nil {
			slog.Error("parse maven dependency txt failed", slog.String("project", repo.Name), slog.String("path", repo.Path), slog.String("err", err.Error()))
			return err
		}
		slog.Info("parse maven dependency txt", slog.Any("project", repo.Name), slog.Any("path", repo.Path), slog.Any("dependencyTxt", dependencyTxt))
		projectPom, err := pom.ParsePOM(txtPath)
		if err != nil {
			slog.Error("parse maven pom failed", slog.String("project", repo.Name), slog.String("path", repo.Path), slog.String("err", err.Error()))
			return err
		}
		slog.Info("parse maven pom", slog.Any("project", repo.Name), slog.Any("path", repo.Path), slog.Any("pom", projectPom))
	}
	return nil
}
