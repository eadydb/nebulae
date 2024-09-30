package nebulae

import (
	"context"

	"github.com/spf13/cobra"
)

func NewCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Args:          cobra.NoArgs,
		Use:           "nebulae [command]",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand()
	return cmd
}
