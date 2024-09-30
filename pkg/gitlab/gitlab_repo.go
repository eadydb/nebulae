package gitlab

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/eadydb/nebulae/pkg/repository"
	"gorm.io/gorm"
)

type IRepository interface {
	SaveRepository(repository Repository) error
	UpdateRepository(language, dir string, id int) error
	FindRepositories(page, limit int, language string) ([]GitlabRepository, error)
	GetRepositoryById(id int) (*GitlabRepository, error)
	ExixtsRepository(id []int) ([]int, error)
}

type RepositoryService struct {
	Db  *gorm.DB
	Ctx context.Context
}

func newRepositoryService(ctx context.Context) *RepositoryService {
	return &RepositoryService{
		Db:  repository.GormConn,
		Ctx: ctx,
	}
}

func (r *RepositoryService) SaveRepository(repository Repository) error {
	sql := fmt.Sprintf("INSERT INTO t_gitlab_repository (id, name, path, path_with_namespace, default_branch, archived, http_url_to_repo, web_url, last_activity_at, created_at, visibility, description, creator_id, namespace_id, namespace_name, namespace_path, namespace_kind, namespace_full_path, namespace_parent_id) VALUES (%d, '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, '%s', '%s', '%s', '%s', %d)",
		repository.Id, repository.Name, repository.Path, repository.PathWithNamespace, repository.DefaultBranch, repository.Archived, repository.HttpUrlToRepo, repository.WebUrl, repository.LastActivityAt, repository.CreateAt, repository.Visibility, repository.Description, repository.CreatorId, repository.Namespace.Id, repository.Namespace.Name, repository.Namespace.Path, repository.Namespace.Kind, repository.Namespace.FullPath, repository.Namespace.ParentId)
	err := r.Db.Exec(sql).Error
	if err != nil {
		slog.Error("save gitlab repository failed", slog.String("err", err.Error()), slog.Any("repository", repository))
		return err
	}
	return nil
}

func (r *RepositoryService) UpdateRepository(language, dir string, id int) error {
	err := r.Db.Exec(fmt.Sprintf("update t_gitlab_repository set language = '%s', workspace_dir = '%s' where id = %d", language, dir, id)).Error
	if err != nil {
		slog.Error("update gitlab repository language failed", slog.String("err", err.Error()), slog.Int("id", id), slog.String("language", language))
		return err
	}
	return nil
}

func (r *RepositoryService) FindRepositories(page, limit int, language string) ([]GitlabRepository, error) {
	var repositories []GitlabRepository
	sql := fmt.Sprintf("select * from t_gitlab_repository limit %d offset %d order by id", limit, (page-1)*limit)
	if language != "" {
		sql = fmt.Sprintf("select * from t_gitlab_repository where language match_phrase '%s' limit %d offset %d order by id", language, limit, (page-1)*limit)
	}
	err := r.Db.Raw(sql).Scan(&repositories).Error
	if err != nil {
		slog.Error("find gitlab repositories failed", slog.String("err", err.Error()), slog.Int("page", page), slog.Int("limit", limit), slog.String("language", language))
		return nil, err
	}
	return repositories, nil
}

func (r *RepositoryService) GetRepositoryById(id int) (*GitlabRepository, error) {
	var repository GitlabRepository
	err := r.Db.Raw(fmt.Sprintf("select * from t_gitlab_repository where id = %d", id)).Scan(&repository).Error
	if err != nil {
		slog.Error("find gitlab repository by id failed", slog.String("err", err.Error()), slog.Int("id", id))
		return nil, err
	}
	return &repository, nil
}

func (r *RepositoryService) ExixtsRepository(ids []string) ([]int, error) {
	var exists []int
	err := r.Db.Raw(fmt.Sprintf("select id from t_gitlab_repository where id in (%s)", strings.Join(ids, ","))).Scan(&exists).Error
	if err != nil {
		slog.Error("find gitlab repository by ids failed", slog.String("err", err.Error()), slog.Any("ids", ids))
		return nil, err
	}
	return exists, nil
}

func saveRepository(ctx context.Context, repositories []Repository) error {
	ids := make([]string, 0)
	for _, repository := range repositories {
		ids = append(ids, strconv.Itoa(repository.Id))
	}
	repoService := newRepositoryService(ctx)
	existsIds, err := repoService.ExixtsRepository(ids)
	if err != nil {
		slog.Error("exists repository failed", slog.String("err", err.Error()))
		return err
	}
	for _, repository := range repositories {
		if !exists(existsIds, repository.Id) {
			if err := repoService.SaveRepository(repository); err != nil {
				slog.Error("save repository failed", slog.String("err", err.Error()))
				return err
			}
		}
	}
	return nil
}

func exists(ids []int, id int) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}
