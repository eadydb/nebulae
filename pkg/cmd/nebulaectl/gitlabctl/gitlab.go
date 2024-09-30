package gitlabctl

import (
	"context"
	"log/slog"

	"github.com/eadydb/nebulae/pkg/config"
	"github.com/eadydb/nebulae/pkg/repository"
	"github.com/spf13/cobra"
)

func NewGitlabCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Args:          cobra.NoArgs,
		Use:           "gitlab [command]",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		ScanGitlabHubRepository(ctx),
		CloneRepository(ctx),
	)
	return cmd
}

func initDriver(ctx context.Context, driver *config.Driver) {
	if driver == nil {
		return
	}
	if driver.Mysql != nil {
		conn, err := repository.NewMysqlDriver("", driver.Mysql.Host, driver.Mysql.Username, driver.Mysql.Password, driver.Mysql.Database, driver.Mysql.Port).GetGromConn(ctx)
		if err != nil {
			slog.Error("get gorm connection failed", slog.String("err", err.Error()))
			panic(err)
		}
		repository.GormConn = conn
	}
	if driver.Neo4j != nil {
		err := repository.NewNeo4jDriver(driver.Neo4j.Host, driver.Neo4j.Username, driver.Neo4j.Password, driver.Neo4j.Port).Connect()
		if err != nil {
			slog.Error("connect neo4j failed", slog.String("err", err.Error()))
			panic(err)
		}
	}
}
