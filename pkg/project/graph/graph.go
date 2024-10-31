package graph

import (
	"context"
	"log/slog"

	"github.com/eadydb/nebulae/pkg/gitlab"
	"github.com/eadydb/nebulae/pkg/project/apollo"
	"github.com/eadydb/nebulae/pkg/project/maven"
	"github.com/eadydb/nebulae/pkg/project/pom"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Graph struct {
	Ctx              context.Context
	Session          neo4j.SessionWithContext
	GitlabRepository gitlab.GitlabRepository
}

func NewGraph(ctx context.Context, gitlabRepository gitlab.GitlabRepository, sesesion neo4j.SessionWithContext) *Graph {
	return &Graph{
		Ctx:     ctx,
		Session: sesesion,
	}
}

func (g *Graph) Do(pom *pom.POM, dependencyTxt *maven.DependencyTxt, middlewares []apollo.Middleware, appId string) error {
	return nil
}

func (g *Graph) doProject() error {
	_, err := g.Session.ExecuteWrite(g.Ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MERGE (project:Project {name: $name, language: $language, description: $description, create_at: $create_at})
			return project
		`
		params := map[string]interface{}{
			"name":        g.GitlabRepository.Name,
			"language":    g.GitlabRepository.Language,
			"description": g.GitlabRepository.Description,
			"create_at":   g.GitlabRepository.CreateAt,
		}
		_, err := tx.Run(g.Ctx, query, params)
		if err != nil {
			slog.Error("create project failed", slog.String("err", err.Error()), slog.String("projectName", g.GitlabRepository.Name))
		}
		return nil, err
	})

	return err
}

func (g *Graph) doProjectGitlab() error {
	_, err := g.Session.ExecuteWrite(g.Ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MERGE (gitlab:Gitlab {id: $id})
			SET gitlab.name = $name,
				gitlab.http_url_to_repo = $http_url_to_repo,
				gitlab.path_with_namespace = $path_with_namespace,
				gitlab.default_branch = $default_branch,
				gitlab.create_at = $create_at,
				gitlab.namespace_full_path = $namespace_full_path
			MERGE (project:Project {name: $projectName})
			MERGE (project)-[:HAS_GITLAB]->(gitlab)
			RETURN gitlab
		`
		params := map[string]interface{}{
			"id":                  g.GitlabRepository.Id,
			"name":                g.GitlabRepository.Name,
			"http_url_to_repo":    g.GitlabRepository.HttpUrlToRepo,
			"path_with_namespace": g.GitlabRepository.PathWithNamespace,
			"default_branch":      g.GitlabRepository.DefaultBranch,
			"create_at":           g.GitlabRepository.CreateAt,
			"namespace_full_path": g.GitlabRepository.NamespaceFullPath,
			"namespace":           g.GitlabRepository.NamespaceName,
			"projectName":         g.GitlabRepository.Name,
		}
		_, err := tx.Run(g.Ctx, query, params)
		if err != nil {
			slog.Error("create gitlab failed", slog.String("err", err.Error()), slog.String("projectName", g.GitlabRepository.Name))
		}
		return nil, err
	})
	return err
}

func (g *Graph) doProjectApollo(middlewares []apollo.Middleware, appId string) error {
	_, err := g.Session.ExecuteWrite(g.Ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MERGE (apollo:Apollo {appId: $appId})
			MERGE (project:Project {name: $projectName})
			MERGE (project)-[:HAS_APOLLO]->(apollo)
			RETURN apollo
		`
		params := map[string]interface{}{
			"appId":       appId,
			"projectName": g.GitlabRepository.Name,
		}
		if _, err := tx.Run(g.Ctx, query, params); err != nil {
			slog.Error("create apollo failed", slog.String("err", err.Error()), slog.String("projectName", g.GitlabRepository.Name))
			return nil, err
		}

		for _, middleware := range middlewares {
			query = `
				MERGE (middleware:Middleware {name: $name})
				SET middleware.type = $type
				MERGE (project:Project {name: $projectName})
				MERGE (project)-[:HAS_MIDDLEWARE]->(middleware)
				RETURN middleware
			`
			params = map[string]interface{}{
				"name":        middleware.Name,
				"type":        middleware.Type,
				"projectName": g.GitlabRepository.Name,
			}
			_, err := tx.Run(g.Ctx, query, params)
			if err != nil {
				slog.Error("create project middleware failed", slog.String("err", err.Error()), slog.String("projectName", g.GitlabRepository.Name))
				return nil, err
			}
		}

		return nil, nil
	})
	return err
}

func (g *Graph) doProjectMaven() error {
	return nil
}

func (g *Graph) doProjectPom(pom *pom.POM) error {
	_, err := g.Session.ExecuteWrite(g.Ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		pomName := pom.Name
		if pomName == "" || len(pomName) == 0 {
			pomName = pom.ArtifactID
		}

		query := `
			MERGE (pom:Pom {name: $name, groupId: $groupId, artifactId: $artifactId})
			MERGE (parent:Parent {group_id: $parentGroupId, artifact_id: $parentArtifactId, version: $parentVersion})
			MERGE (project:Project {name: $projectName})
			MERGE (pom)-[:HAS_PARENT]->(parent)
			MERGE (project)-[:HAS_POM]->(pom)
			RETURN pom
		`
		params := map[string]interface{}{
			"groupId":          pom.GroupID,
			"artifactId":       pom.GroupID,
			"name":             pomName,
			"parentGroupId":    pom.Parent.GroupID,
			"parentArtifactId": pom.Parent.ArtifactID,
			"parentVersion":    pom.Parent.Version,
			"projectName":      g.GitlabRepository.Name,
		}
		if _, err := tx.Run(g.Ctx, query, params); err != nil {
			slog.Error("create pom failed", slog.String("err", err.Error()), slog.String("projectName", g.GitlabRepository.Name))
			return nil, err
		}

		for _, module := range pom.Modules {
			query = `
				MERGE (module:Module {name: $name})
				MERGE (pom:Pom {name: $pomName})
				MERGE (project)-[:HAS_MODULE]->(module)
				RETURN module
			`
			params = map[string]interface{}{
				"name":    module,
				"pomName": pomName,
			}
			_, err := tx.Run(g.Ctx, query, params)
			if err != nil {
				slog.Error("create module failed", slog.String("err", err.Error()), slog.String("projectName", g.GitlabRepository.Name))
				return nil, err
			}
		}
		return nil, nil
	})
	return err
}
