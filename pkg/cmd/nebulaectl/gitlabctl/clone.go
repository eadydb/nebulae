package gitlabctl

import (
	"context"
	"log/slog"

	"github.com/eadydb/nebulae/pkg/config"
	"github.com/eadydb/nebulae/pkg/gitlab"
	"github.com/spf13/cobra"
)

func CloneRepository(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Args:          cobra.NoArgs,
		Use:           "clone [command]",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	return cmd
}

func cloneRunE(ctx context.Context) error {
	options := ctx.Value(0).(config.NebulaeOptions)
	slog.Info("scan gitlab hub repository", slog.String("options", options.Gitlab.GitlabUrl))
	initDriver(ctx, options.Driver)

	projects := &gitlab.CloneProjects{
		Page:         1,
		PerPage:      10,
		Recursive:    true,
		WorkspaceDir: options.Gitlab.WorkspaceDir,
	}

	return gitlab.NewGitlab(ctx, options.Gitlab.PrivateToken, options.Gitlab.GitlabUrl).CloneGitlabHub(projects)
}
