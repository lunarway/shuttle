// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/pkg/errors"
	"fmt"
)

var shell string

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion <shell>",
	Short: "Generates shell completion scripts",
	Long: `Output shell completion code for the specified shell (bash or zsh). The shell code must be evaluated to provide
interactive completion of shuttle commands. This can be done by sourcing it from the .bash _profile.

Detailed instructions on how to do this are available here:
https://kubernetes.io/docs/tasks/tools/install-shuttle/#enabling-shell-autocompletion

Note for zsh users: [1] zsh completions are only supported in versions of zsh >= 5.2

Examples:
  # Installing bash completion on macOS using homebrew
  ## If running Bash 3.2 included with macOS
  brew install bash-completion
  ## or, if running Bash 4.1+
  brew install bash-completion@2
  ## If shuttle is installed via homebrew, this should start working immediately.
  ## If you've installed via other means, you may need add the completion to your completion directory
  shuttle completion bash > $(brew --prefix)/etc/bash_completion.d/shuttle


  # Installing bash completion on Linux
  ## Load the shuttle completion code for bash into the current shell
  source <(shuttle completion bash)
  ## Write bash completion code to a file and source if from .bash_profile
  shuttle completion bash > ~/.kube/completion.bash.inc
  printf "
  # Kubectl shell completion
  source '$HOME/.kube/completion.bash.inc'
  " >> $HOME/.bash_profile
  source $HOME/.bash_profile

  # Load the shuttle completion code for zsh[1] into the current shell
  source <(shuttle completion zsh)
  # Set the shuttle completion code for zsh[1] to autoload on startup
  shuttle completion zsh > "${fpath[1]}/_kubectl"`,
	ValidArgs: []string{"bash", "zsh"},
	Args: func(cmd *cobra.Command, args []string) error {
		if cobra.ExactArgs(1)(cmd, args) != nil || cobra.OnlyValidArgs(cmd, args) != nil {
			return errors.New(fmt.Sprintf("Only %v arguments are allowed", cmd.ValidArgs))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
 		switch args[0] {
		case "zsh":
			rootCmd.GenZshCompletion(os.Stdout)
		case "bash":
			rootCmd.GenBashCompletion(os.Stdout)
		default:
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
