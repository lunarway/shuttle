package cmd

import (
	"github.com/lunarway/shuttle/internal/extensions"
	"github.com/spf13/cobra"
)

func newExtCmd() *cobra.Command {
	extManager := extensions.NewExtensionsManager("some registry")

	cmd := &cobra.Command{
		Use:  "ext",
		Long: "helps you manage shuttle extensions",
	}

	cmd.AddCommand(
		newExtInstallCmd(extManager),
		newExtUpdateCmd(extManager),
		newExtInitCmd(extManager),
	)

	return cmd
}

func newExtInstallCmd(extManager *extensions.ExtensionsManager) *cobra.Command {
	cmd := &cobra.Command{
		Use: "install",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func newExtUpdateCmd(extManager *extensions.ExtensionsManager) *cobra.Command {
	cmd := &cobra.Command{
		Use: "update",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func newExtInitCmd(extManager *extensions.ExtensionsManager) *cobra.Command {
	cmd := &cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
