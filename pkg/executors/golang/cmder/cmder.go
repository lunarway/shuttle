package cmder

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/lunarway/shuttle/pkg/executors/golang/executer"
	"github.com/spf13/cobra"
)

type RootCmd struct {
	Cmds []*Cmd
}

func NewRoot() *RootCmd {
	cmd := &RootCmd{}

	return cmd
}

func (rc *RootCmd) AddCmds(cmd ...*Cmd) *RootCmd {
	rc.Cmds = append(rc.Cmds, cmd...)

	return rc
}

func (rc *RootCmd) Execute() {
	if err := rc.TryExecute(os.Args[1:]); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func (rc *RootCmd) TryExecute(args []string) error {
	rootcmd := &cobra.Command{Use: "actions"}

	rootcmd.AddCommand(
		&cobra.Command{Use: "ls", RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Parent().Help()
		}},
	)

	rootcmd.AddCommand(
		&cobra.Command{
			Hidden: true,
			Use:    "lsjson",
			RunE: func(cmd *cobra.Command, args []string) error {
				actions := executer.NewActions()
				for _, cmd := range rc.Cmds {
					args := make([]executer.ActionArg, 0)

					for _, arg := range cmd.Args {
						args = append(args, executer.ActionArg{
							Name: arg.Name,
						})
					}

					actions.Actions[cmd.Name] = executer.Action{
						Args: args,
					}
				}

				rawJson, err := json.Marshal(actions)
				if err != nil {
					return err
				}

				// Prints the commands and args as json to stdout
				_, err = fmt.Printf("%s", string(rawJson))
				if err != nil {
					return err
				}

				return nil
			},
		},
	)

	for _, cmd := range rc.Cmds {
		cmd := cmd
		parameters := make([]string, len(cmd.Args))

		cobracmd := &cobra.Command{
			Use: cmd.Name,
			RunE: func(cobracmd *cobra.Command, args []string) error {
				if err := cobracmd.ParseFlags(args); err != nil {
					return err
				}

				inputs := make([]reflect.Value, 0, len(cmd.Args)+1)
				inputs = append(inputs, reflect.ValueOf(context.Background()))
				for _, arg := range parameters {
					inputs = append(inputs, reflect.ValueOf(arg))
				}

				returnValues := reflect.
					ValueOf(cmd.Func).
					Call(inputs)

				if len(returnValues) == 0 {
					return nil
				}

				for _, val := range returnValues {
					if val.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
						err, ok := val.Interface().(error)
						if ok && err != nil {
							return err
						}
					}
				}

				return nil
			},
		}
		for i, arg := range cmd.Args {
			cobracmd.PersistentFlags().StringVar(&parameters[i], arg.Name, "", "")
			_ = cobracmd.MarkPersistentFlagRequired(arg.Name)
		}

		rootcmd.AddCommand(cobracmd)
	}

	rootcmd.SetArgs(args)
	if err := rootcmd.Execute(); err != nil {
		return err
	}
	return nil
}

type Arg struct {
	Name string
}

type Cmd struct {
	Name string
	Func any
	Args []Arg
}

func NewCmd(name string, f any) *Cmd {
	return &Cmd{
		Name: name,
		Func: f,
		Args: []Arg{},
	}
}

func WithArgs(cmd *Cmd, argName string) *Cmd {
	cmd.Args = append(cmd.Args, Arg{Name: argName})
	return cmd
}
