package gitlab

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/eadydb/nebulae/pkg/utils/network"
	"github.com/go-git/go-git/config"
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
	repositories, err := g.SelectGitlabRepository(projects)
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
	if len(repositories) < projects.PerPage {
		return nil
	}
	projects.Page++
	return g.ScanGitlabHub(projects)
}

// SelectGitlabRepository select gitlab repository
func (g *Gitlab) SelectGitlabRepository(projects *GitlabProjects) ([]Repository, error) {
	params := make(map[string]string)
	params["private_token"] = g.PrivateToken
	params["simple"] = strconv.FormatBool(projects.Simple)
	params["membership"] = strconv.FormatBool(projects.Membership)
	params["page"] = strconv.Itoa(projects.Page)
	params["per_page"] = strconv.Itoa(projects.PerPage)

	var repositories []Repository
	return network.Get(g.Domain, "/api/v4/projects", "", params, repositories)
}

// Clone clone gitlab repository
func (g *Gitlab) Clone(repository GitlabRepository, hubDir string) (string, error) {
	dir := hubDir + "/" + repository.PathWithNamespace
	var repo *git.Repository
	var err error
	if existsLocalRepository(dir) {
		repo, err = git.PlainOpen(dir)
		if err != nil {
			return "", err
		}
	} else {
		url := strings.ReplaceAll(repository.HttpUrlToRepo, "http://", "https://")
		repo, err = git.PlainClone(dir, false, &git.CloneOptions{
			URL: url,
			Auth: &http.BasicAuth{
				Username: g.Username,
				Password: g.PrivateToken,
			},
		})
		if err != nil {
			slog.Error("Gitlab Repository failed", slog.String("err", err.Error()), slog.String("url", repository.HttpUrlToRepo))
			return "", err
		}
	}
	wt, err := repo.Worktree()
	if err != nil {
		slog.Error("Get Gitlab Repository Worktree failed", slog.String("err", err.Error()), slog.String("url", repository.HttpUrlToRepo))
		return "", err
	}
	ref, err := repo.Head()
	if err != nil {
		slog.Error("Get Gitlab Repository HEAD failed", slog.String("err", err.Error()), slog.String("url", repository.HttpUrlToRepo))
		return "", err
	}

	create := ref.Name().Short() != repository.DefaultBranch
	if create {
		branchRefName := plumbing.NewBranchReferenceName(repository.DefaultBranch)
		branchCoOpts := git.CheckoutOptions{
			Branch: plumbing.ReferenceName(branchRefName),
			Force:  true,
		}
		if err := wt.Checkout(&branchCoOpts); err != nil {
			mirrorRemoteBranchRefSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", repository.DefaultBranch, repository.DefaultBranch)
			if err = g.fetchOrigin(repo, mirrorRemoteBranchRefSpec); err != nil {
				slog.Error("Fetch Remote Branch failed", slog.String("err", err.Error()), slog.String("url", repository.HttpUrlToRepo), slog.String("branch", repository.DefaultBranch))
				return dir, err
			}

			if err = wt.Checkout(&branchCoOpts); err != nil {
				slog.Error("Checkout Remote Branch failed", slog.String("err", err.Error()), slog.String("url", repository.HttpUrlToRepo), slog.String("branch", repository.DefaultBranch))
				return dir, err
			}

		}
		slog.Info("Checkout Remote Branch success", slog.String("branch", repository.DefaultBranch), slog.String("url", repository.HttpUrlToRepo))
	}

	return dir, nil
}

func (g *Gitlab) fetchOrigin(repo *git.Repository, refSpecStr string) error {
	remote, err := repo.Remote("origin")
	if err != nil {
		slog.Error("获取远程分支Origin失败", slog.String("err", err.Error()))
		return err
	}

	var refSpecs []config.RefSpec
	if refSpecStr != "" {
		refSpecs = []config.RefSpec{config.RefSpec(refSpecStr)}
	}

	if err = remote.Fetch(&git.FetchOptions{
		RefSpecs: refSpecs,
		Auth: &http.BasicAuth{
			Username: g.Username,
			Password: g.PrivateToken,
		},
	}); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			slog.Info("refs already up to date")
		} else {
			return fmt.Errorf("fetch origin failed: %v", err)
		}
	}

	return nil
}

// existsLocalRepository check if the repository exists
func existsLocalRepository(repository string) bool {
	if _, err := os.Stat(repository); err == nil {
		return true
	}
	return false
}
