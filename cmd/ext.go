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
		Use:  "install",
		Long: "Install ensures that extensions already known about are downloaded and available",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func newExtUpdateCmd(extManager *extensions.ExtensionsManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "update",
		Long: "Update will fetch the latest version of the extensions from the given registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func newExtInitCmd(extManager *extensions.ExtensionsManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "init",
		Long: "init will create an initial extensions repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
