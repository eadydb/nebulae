package config

import "github.com/eadydb/nebulae/pkg/repository"

type NebulaeOptions struct {
	EnableRPC   bool
	RPCPort     int
	RPCHTTPPort int
	Driver      *Driver
	Gitlab      *GitlabConfig
}

type Driver struct {
	Mysql *repository.MysqlDriver
	Neo4j *repository.Neo4jDriver
}

type GitlabConfig struct {
	PrivateToken   string
	GitlabUrl      string
	GitlabUsername string
	WorkspaceDir   string
}
