package gitlabctl

import (
	"context"
	"log/slog"

	"github.com/eadydb/nebulae/pkg/config"
	"github.com/eadydb/nebulae/pkg/gitlab"
	"github.com/spf13/cobra"
)

func ScanGitlabHubRepository(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Args:          cobra.NoArgs,
		Use:           "scan [command]",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return scanRunE(ctx)
		},
	}
	return cmd
}

func scanRunE(ctx context.Context) error {
	options := ctx.Value(0).(config.NebulaeOptions)
	slog.Info("scan gitlab hub repository", slog.String("options", options.Gitlab.GitlabUrl))
	initDriver(ctx, options.Driver)
	projects := &gitlab.GitlabProjects{
		Simple:     true,
		Page:       1,
		PerPage:    100,
		Recursive:  true,
		Membership: false,
	}
	return gitlab.NewGitlab(ctx, options.Gitlab.PrivateToken, options.Gitlab.GitlabUrl).ScanGitlabHub(projects)
}
