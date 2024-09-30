package nebulaectl

import (
	"context"

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
	cmd.AddCommand()
	return cmd
}
