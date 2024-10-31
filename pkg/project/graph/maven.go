package graph

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/eadydb/nebulae/pkg/gitlab"
	"github.com/eadydb/nebulae/pkg/project/apollo"
	"github.com/eadydb/nebulae/pkg/project/maven"
	"github.com/eadydb/nebulae/pkg/project/pom"
	"github.com/eadydb/nebulae/pkg/repository"
	"github.com/eadydb/nebulae/pkg/utils/walk"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"gorm.io/gorm"
)

type GitlabRepositoryMaven struct {
	Ctx      context.Context
	DB       *gorm.DB
	Driver   neo4j.DriverWithContext
	PageSize int
	Language string
}

func NewGitlabRepositoryMaven(ctx context.Context, pageSize int, language string) *GitlabRepositoryMaven {
	return &GitlabRepositoryMaven{
		PageSize: pageSize,
		Ctx:      ctx,
		DB:       repository.GormConn,
		Language: language,
		Driver:   repository.Neo4jConn,
	}
}

func (g *GitlabRepositoryMaven) Do() error {
	session := g.Driver.NewSession(g.Ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(g.Ctx)
	return nil
}

func (g *GitlabRepositoryMaven) DoMvnCommand(repo *gitlab.RepositoryService, limit int, session neo4j.SessionWithContext) error {
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
			err := g.DoMvnDependency(repo, session)
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
		return g.DoMvnCommand(repo, limit+1, session)
	}
	return nil
}

func (g *GitlabRepositoryMaven) DoMvnDependency(repo gitlab.GitlabRepository, session neo4j.SessionWithContext) error {
	mvn := maven.NewMaven(g.Ctx, repo.WorkspaceDir, repo.Name)
	if err := mvn.ExectueCmd(true); err != nil {
		slog.Error("execute maven dependency:tree failed", slog.String("project", repo.Name), slog.String("path", repo.Path), slog.String("err", err.Error()))
		return err
	}
	txtPaths, err := mvn.LoadingMvnTxtFile()
	if err != nil {
		slog.Error("loading maven txt file failed", slog.String("project", repo.Name), slog.String("path", repo.Path), slog.String("err", err.Error()))
		return err
	}

	var wg sync.WaitGroup
	maxConcurrent := 2
	sem := make(chan struct{}, maxConcurrent)
	errChan := make(chan error, len(txtPaths))

	for _, txtPath := range txtPaths {
		sem <- struct{}{}
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("panic in goroutine: %v", r)
				}
				<-sem
			}()
			err := g.DoParse(path, repo, session)
			if err != nil {
				errChan <- err
			}
		}(txtPath)

	}
	return nil
}

func (g *GitlabRepositoryMaven) DoParse(txtPath string, repo gitlab.GitlabRepository, session neo4j.SessionWithContext) error {
	dependencyTxt, err := maven.ParseMavenDependencyTxt(txtPath)
	if err != nil {
		slog.Error("parse maven dependency txt failed", slog.String("project", repo.Name), slog.String("path", repo.Path), slog.String("err", err.Error()))
		return err
	}
	pomXml := txtPath[:len(txtPath)-len(maven.Suffix)] + "pom.xml"
	projectPom, err := pom.ParsePOM(pomXml)
	if err != nil {
		slog.Error("parse maven pom failed", slog.String("project", repo.Name), slog.String("path", repo.Path), slog.String("err", err.Error()))
		return err
	}

	// load project boostrap or application
	ymlPaths, err := walk.From(repo.Path).CollectFilterPaths("application.yml")
	if len(ymlPaths) == 0 {
		ymlPaths, _ = walk.From(repo.Path).CollectFilterPaths("bootstrap.yml")
	}
	middlewares := make([]apollo.Middleware, 0)
	appId := ""
	if len(ymlPaths) > 0 {
		yamPath := ymlPaths[0]
		for _, path := range ymlPaths {
			if strings.Contains(path, "/prod/") {
				yamPath = path
			}
		}
		middlewares, _, _ = apollo.UnmarshalApollo(yamPath)
	}

	return NewGraph(g.Ctx, repo, session).Do(projectPom, dependencyTxt, middlewares, appId)
}

// LoadingFile 获取pom.xml文件
func LoadingFile(path, fileName string) ([]string, error) {
	txtPaths, err := walk.From(path).CollectFilterPaths(fileName)
	if err != nil {
		slog.Error("loading mvn txt file failed", slog.String("dir", path), slog.String("err", err.Error()), slog.String("fileName", fileName))
		return nil, err
	}
	return txtPaths, nil
}
