package plugins

import (
	"github.com/lunarway/shuttle/cmd/utility"
	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/spf13/cobra"
)

func NewPluginsRootCmd(uii *ui.UI, contextProvider utility.ContextProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugins [command]",
		Short: "Interact with plugins",
	}

	cmd.AddCommand(
		newInitCmd(uii, contextProvider),
		newLsCmd(uii, contextProvider),
	)

	return cmd
}
