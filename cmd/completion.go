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
	"bytes"
	"fmt"
	"io"

	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/spf13/cobra"
)

func newCompletion(uii *ui.UI) *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion <shell>",
		Short: `Output shell completion code`,
		Long: `Output shell completion code for the specified shell (bash or zsh).
The shell code must be evaluated to provide interactive
completion of shuttle commands.  This can be done by sourcing it from
the .bash_profile.

Note for zsh users: zsh completions are only supported in versions of zsh >= 5.2

Installing bash completion on macOS using homebrew

    If running Bash 3.2 included with macOS

    	brew install bash-completion

    If running Bash 4.1+

    	brew install bash-completion@2

    You may need add the completion to your completion directory

    	shuttle completion bash > $(brew --prefix)/etc/bash_completion.d/shuttle

Installing bash completion on Linux

    If bash-completion is not installed on Linux, please install the 'bash-completion' package
    via your distribution's package manager.

    Load the shuttle completion code for bash into the current shell

    	source <(shuttle completion bash)

    Write bash completion code to a file and source if from .bash_profile

     	shuttle completion bash > ~/.shuttle/completion.bash.inc
     	printf "
     	            # shuttle shell completion
     	source '$HOME/.shuttle/completion.bash.inc'
     	            " >> $HOME/.bash_profile
    	source $HOME/.bash_profile

    Load the shuttle completion code for zsh[1] into the current shell

    	source <(shuttle completion zsh)

    Set the shuttle completion code for zsh[1] to autoload on startup

    	shuttle completion zsh > "${fpath[1]}/_shuttle"`,
		ValidArgs: []string{"bash", "zsh"},
		Args: func(cmd *cobra.Command, args []string) error {
			if cobra.ExactArgs(1)(cmd, args) != nil || cobra.OnlyValidArgs(cmd, args) != nil {
				return fmt.Errorf("only %v arguments are allowed", cmd.ValidArgs)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			*uii = uii.SetContext(ui.LevelSilent)
			switch args[0] {
			case "zsh":
				runCompletionZsh(cmd.OutOrStdout(), cmd.Root())
			case "bash":
				cmd.Root().GenBashCompletion(cmd.OutOrStdout())
			default:
			}
		},
	}

	return completionCmd
}

// this writes a zsh completion script that wraps the bash completion script.
//
// Copied from kubectl: https://github.com/kubernetes/kubernetes/blob/9c2df998af9eb565f11d42725dc77e9266483ffc/pkg/kubectl/cmd/completion/completion.go#L145
func runCompletionZsh(out io.Writer, shuttle *cobra.Command) error {
	zshHead := "#compdef shuttle\n"

	out.Write([]byte(zshHead))

	zshInitialization := `
__shuttle_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand
	source "$@"
}
__shuttle_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__shuttle_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}
__shuttle_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}
__shuttle_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}
__shuttle_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}
__shuttle_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}
__shuttle_filedir() {
	local RET OLD_IFS w qw
	__shuttle_debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi
	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"
	IFS="," __shuttle_debug "RET=${RET[@]} len=${#RET[@]}"
	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__shuttle_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}
__shuttle_quote() {
    if [[ $1 == \'* || $1 == \"* ]]; then
        # Leave out first character
        printf %q "${1:1}"
    else
	printf %q "$1"
    fi
}
autoload -U +X bashcompinit && bashcompinit
# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi
__shuttle_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__shuttle_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__shuttle_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__shuttle_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__shuttle_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__shuttle_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__shuttle_type/g" \
	<<'BASH_COMPLETION_EOF'
`
	out.Write([]byte(zshInitialization))

	buf := new(bytes.Buffer)
	shuttle.GenBashCompletion(buf)
	out.Write(buf.Bytes())

	zshTail := `
BASH_COMPLETION_EOF
}
__shuttle_bash_source <(__shuttle_convert_bash_to_zsh)
_complete shuttle 2>/dev/null
`
	out.Write([]byte(zshTail))
	return nil
}
