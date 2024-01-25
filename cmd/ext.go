package cmd

import (
	"errors"
	"os"
	"os/exec"

	stdcontext "context"

	"github.com/lunarway/shuttle/internal/extensions"
	"github.com/lunarway/shuttle/internal/global"
	"github.com/spf13/cobra"
)

type extGlobalConfig struct {
	registry string
}

func (c *extGlobalConfig) getRegistry() (string, bool) {
	if c.registry != "" {
		return c.registry, true
	}

	if registryEnv := os.Getenv("SHUTTLE_EXTENSIONS_REGISTRY"); registryEnv != "" {
		return registryEnv, true
	}

	return "", false
}

func addExtensions(rootCmd *cobra.Command) error {
	extManager := extensions.NewExtensionsManager(global.NewGlobalStore())

	extensions, err := extManager.GetAll(stdcontext.Background())
	if err != nil {
		return err
	}
	grp := &cobra.Group{
		ID:    "extensions",
		Title: "Extensions",
	}
	rootCmd.AddGroup(grp)
	for _, extension := range extensions {
		extension := extension

		rootCmd.AddCommand(
			&cobra.Command{
				Use:                extension.Name(),
				Short:              extension.Description(),
				Version:            extension.Version(),
				GroupID:            "extensions",
				DisableFlagParsing: true,
				RunE: func(cmd *cobra.Command, args []string) error {

					extCmd := exec.CommandContext(cmd.Context(), extension.FullPath(), args...)

					extCmd.Stdout = os.Stdout
					extCmd.Stderr = os.Stderr
					extCmd.Stdin = os.Stdin

					if err := extCmd.Start(); err != nil {
						return err
					}

					if err := extCmd.Wait(); err != nil {
						return err
					}

					return nil
				},
			},
		)
	}

	return nil
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

			registry, ok := globalConfig.getRegistry()
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
