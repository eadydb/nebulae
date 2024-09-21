package apollo

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	middlewarePrefix = []string{"kf-", "-redis", "ds-", "rocket-", "sentinel-", "rabbit-", "es-", "ds-sb-", "ds-selectdb-", "durid-", "hikari-"}
	middlewareMap    = map[string]string{
		"kf-":          "Kafka",
		"-redis":       "Redis",
		"ds-":          "DataSource",
		"druid-":       "DataSource",
		"hikari-":      "DataSource",
		"rocket-":      "RocketMQ",
		"sentinel-":    "Sentinel",
		"rabbit-":      "RabbitMQ",
		"es-":          "ElasticSearch",
		"ds-sb-":       "SelectDB",
		"ds-selectdb-": "SelectDB",
	}
)

// Bootstrap represents application bootstrap.yml config
type Bootstrap struct {
	Apollo Apollo `yaml:"apollo"`
	App    App    `yaml:"app"`
}

// App represents apollo app id
type App struct {
	Id string `yaml:"id"`
}

// Apollo represents apollo config
type Apollo struct {
	Bootstrap ApolloBootstrap `yaml:"bootstrap"`
	Meta      string          `yaml:"meta"`
}

// ApolloBootstrap represents apollo bootstrap config
type ApolloBootstrap struct {
	Enabled    bool   `yaml:"enabled"`
	Namespaces string `yaml:"namespaces"`
}

// Middleware represents application middleware
type Middleware struct {
	Name string // namespace name
	Type string // middleware type
}

// UnmarshalApollo unmarshal apollo config
func UnmarshalApollo(fileName string) ([]Middleware, string, error) {
	if fileName == "" {
		return nil, "", errors.New("fileName is empty")
	}
	slog.Info("unmarshal apollo config", slog.String("filename", fileName))
	content, err := os.ReadFile(fileName)
	if err != nil {
		slog.Error("read apollo config failed", slog.String("error", err.Error()))
		return nil, "", err
	}
	return UnmarshalApolloContent(content)
}

// UnmarshalApolloContent unmarshal apollo content
func UnmarshalApolloContent(content []byte) ([]Middleware, string, error) {
	if len(content) == 0 {
		return nil, "", errors.New("content is empty")
	}
	slog.Info("unmarshal apollo content", slog.String("content", string(content)))
	bootstrap := Bootstrap{}
	err := yaml.Unmarshal(content, &bootstrap)
	if err != nil {
		slog.Error("unmarshal apollo content failed", slog.String("error", err.Error()))
		return nil, "", err
	}
	namespaces, err := bootstrap.Apollo.GetMiddleware()
	return namespaces, bootstrap.App.Id, err
}

// GetMiddleware get middleware from apollo config
func (a *Apollo) GetMiddleware() ([]Middleware, error) {
	if a.Bootstrap.Namespaces == "" {
		return nil, errors.New("namespaces is empty")
	}
	middlewares := make([]Middleware, 0)
	for _, namespace := range strings.Split(a.Bootstrap.Namespaces, ",") {
		if ns, err := parseMiddleware(namespace); err == nil {
			middlewares = append(middlewares, *ns)
		}
	}
	return middlewares, nil
}

func (a *Apollo) GetMeta() string {
	return a.Meta
}

// parseMiddleware parse middleware from namespace
func parseMiddleware(namespace string) (*Middleware, error) {
	if namespace == "" {
		return nil, errors.New("namespace is empty")
	}
	ns := strings.Split(namespace, ".")[0]
	if ok, prefix := existsPrefix(ns); ok {
		middleware := strings.ReplaceAll(ns, prefix, "")
		return &Middleware{Name: middleware, Type: middlewareMap[prefix]}, nil
	}
	return nil, errors.New("middleware not found")
}

// existsPrefix check if namespace contains middleware prefix
func existsPrefix(namespace string) (bool, string) {
	for _, prefix := range middlewarePrefix {
		if strings.Contains(namespace, prefix) {
			return true, prefix
		}
	}
	return false, ""
}
