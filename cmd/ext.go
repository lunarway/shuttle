package cmd

import (
	"errors"
	"os"

	"github.com/lunarway/shuttle/internal/extensions"
	"github.com/lunarway/shuttle/internal/global"
	"github.com/spf13/cobra"
)

type extGlobalConfig struct {
	registry string
}

func (c *extGlobalConfig) Registry() (string, bool) {
	if c.registry != "" {
		return c.registry, true
	}

	if registryEnv := os.Getenv("SHUTTLE_EXTENSIONS_REGISTRY"); registryEnv != "" {
		return registryEnv, true
	}

	return "", false
}

func newExtCmd() *cobra.Command {
	globalConfig := &extGlobalConfig{}

	cmd := &cobra.Command{
		Use:  "ext",
		Long: "helps you manage shuttle extensions",
	}

	cmd.AddCommand(
		newExtInstallCmd(globalConfig),
		newExtUpdateCmd(globalConfig),
		newExtInitCmd(globalConfig),
	)

	cmd.PersistentFlags().StringVar(&globalConfig.registry, "registry", "", "the given registry, if not set will default to SHUTTLE_EXTENSIONS_REGISTRY")

	return cmd
}

func newExtInstallCmd(globalConfig *extGlobalConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "install",
		Long: "Install ensures that extensions are downloaded and available",
		RunE: func(cmd *cobra.Command, args []string) error {
			extManager := extensions.NewExtensionsManager(global.NewGlobalStore())

			if err := extManager.Install(cmd.Context()); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func newExtUpdateCmd(globalConfig *extGlobalConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update will fetch the latest version of the extensions from the given registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			extManager := extensions.NewExtensionsManager(global.NewGlobalStore())

			registry, ok := globalConfig.Registry()
			if !ok {
				return errors.New("registry is not set")
			}

			if err := extManager.Update(cmd.Context(), registry); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func newExtInitCmd(globalConfig *extGlobalConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init will create an initial extensions repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
