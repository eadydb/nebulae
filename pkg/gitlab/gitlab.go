package gitlab

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/eadydb/nebulae/pkg/utils/network"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type Gitlab struct {
	PrivateToken string
	Username     string
	Domain       string
	Ctx          context.Context
}

type GitlabProjects struct {
	Simple     bool
	Membership bool
	Page       int
	PerPage    int
	Recursive  bool
}

type CloneProjects struct {
	Page         int
	PerPage      int
	Recursive    bool
	Language     string
	WorkspaceDir string
}

func NewGitlab(ctx context.Context, privateToken, domain string) *Gitlab {
	return &Gitlab{
		PrivateToken: privateToken,
		Username:     "dengbin@gmail.com",
		Domain:       domain,
		Ctx:          ctx,
	}
}

// ScanGitlabHub scan gitlab hub
func (g *Gitlab) ScanGitlabHub(projects *GitlabProjects) error {
	repositories, err := g.GetGitlabRepository(projects)
	if err != nil {
		slog.Error("select gitlab repository failed", slog.String("err", err.Error()))
		return err
	}
	if len(repositories) == 0 {
		slog.Info("select gitlab repository not exists", slog.Any("projects", projects))
		return nil
	}
	if err := saveRepository(g.Ctx, repositories); err != nil {
		slog.Error("save gitlab repository failed", slog.String("err", err.Error()))
		return err
	}
	if !projects.Recursive {
		slog.Info("End Sweep gitlab repository", slog.Any("projects", projects))
		return nil
	}

	if len(repositories) < projects.PerPage {
		return nil
	}
	projects.Page++
	return g.ScanGitlabHub(projects)
}

// GetGitlabRepository select gitlab repository
func (g *Gitlab) GetGitlabRepository(projects *GitlabProjects) ([]Repository, error) {
	params := make(map[string]string)
	params["private_token"] = g.PrivateToken
	params["simple"] = strconv.FormatBool(projects.Simple)
	params["membership"] = strconv.FormatBool(projects.Membership)
	params["page"] = strconv.Itoa(projects.Page)
	params["per_page"] = strconv.Itoa(projects.PerPage)

	var repositories []Repository
	return network.Get(g.Domain, "/api/v4/projects", "", params, repositories)
}

func (g *Gitlab) CloneGitlabHub(projects *CloneProjects) error {
	repositories, err := newRepositoryService(g.Ctx).FindRepositories(projects.Page, projects.PerPage, projects.Language)
	if err != nil {
		slog.Error("find repositories failed", slog.String("err", err.Error()))
		return err
	}
	if len(repositories) == 0 {
		slog.Info("find repositories not exists", slog.Any("projects", projects))
		return nil
	}
	var wg sync.WaitGroup
	maxConcurrent := 2
	sem := make(chan struct{}, maxConcurrent)

	for _, repository := range repositories {
		wg.Add(1)
		sem <- struct{}{}
		go func(repo GitlabRepository) {
			g.Clone(&repo, projects.WorkspaceDir)
			<-sem
		}(repository)
	}
	wg.Wait()

	if len(repositories) < projects.PerPage {
		return nil
	}
	projects.Page++
	return g.CloneGitlabHub(projects)
}

// Clone clone gitlab repository
func (g *Gitlab) Clone(repository *GitlabRepository, hubDir string) error {
	directory := hubDir + "/" + repository.PathWithNamespace

	// Check if the directory already exists and contains a git repository
	if _, err := os.Stat(directory); err == nil {
		// Directory exists, try to open the repository
		repo, err := git.PlainOpen(directory)
		if err == nil {
			// Repository exists, perform a pull
			if err := g.PullRepository(repo, repository.DefaultBranch, repository.Name); err != nil {
				return err
			}
			return g.UpdateRepositoryLanguage(repository.Id, directory)
		}
	}

	// Directory doesn't exist or is not a git repository, perform a clone
	url := strings.ReplaceAll(repository.HttpUrlToRepo, "http://", "https://")
	if err := g.CloneRepository(url, repository.DefaultBranch, directory); err != nil {
		return err
	}

	return g.UpdateRepositoryLanguage(repository.Id, directory)
}

func (g *Gitlab) CloneRepository(url, branch, directory string) error {
	// Clone options
	cloneOptions := &git.CloneOptions{
		URL:           url,
		Progress:      os.Stdout,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Auth: &http.BasicAuth{
			Username: g.Username,
			Password: g.PrivateToken,
		},
	}

	// Perform the clone
	_, err := git.PlainClone(directory, false, cloneOptions)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	slog.Info("Cloned Successfully ", slog.String("branch", branch), slog.String("url", url), slog.String("directory", directory))
	return nil
}

func (g *Gitlab) PullRepository(repo *git.Repository, branch, name string) error {
	// Get the working directory for the repository
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Pull the latest changes from the remote
	err = w.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		Progress:      os.Stdout,
		Auth: &http.BasicAuth{
			Username: g.Username,
			Password: g.PrivateToken,
		},
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			slog.Info("Repository is already up to date", slog.String("branch", branch), slog.String("repository", name))
			return nil
		}
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}

	slog.Info("Successfully pulled latest changes", slog.String("branch", branch), slog.String("repository", name))
	return nil
}
