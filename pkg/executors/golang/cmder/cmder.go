package cmder

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

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
	rootcmd := &cobra.Command{Use: "shuttletask"}

	rootcmd.AddCommand(
		&cobra.Command{Use: "ls", RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Parent().Help()
		}},
	)

	rootcmd.AddCommand(&cobra.Command{Hidden: true, Use: "lsjson", RunE: func(cmd *cobra.Command, args []string) error {
		cmdNames := make([]string, len(rc.Cmds))
		for i, cmd := range rc.Cmds {
			cmd := cmd
			cmdNames[i] = cmd.Name
		}

		rawJson, err := json.Marshal(cmdNames)
		if err != nil {
			return err
		}

		_, err = fmt.Printf("%s", string(rawJson))
		if err != nil {
			return err
		}

		return nil
	}})

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

				reflect.
					ValueOf(cmd.Func).
					Call(inputs)
				return nil
			},
		}
		for i, arg := range cmd.Args {
			cobracmd.PersistentFlags().StringVar(&parameters[i], arg.Name, "", "")
			_ = cobracmd.MarkPersistentFlagRequired(arg.Name)
		}

		rootcmd.AddCommand(cobracmd)
	}

	if err := rootcmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
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
