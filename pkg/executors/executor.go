package executors

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/errors"
	"github.com/lunarway/shuttle/pkg/ui"
)

type Registry struct {
	executors []Matcher
}

type Matcher func(config.ShuttleAction) (Executor, bool)
type Executor func(context.Context, *ui.UI, ActionExecutionContext) error

func NewRegistry(executors ...Matcher) *Registry {
	return &Registry{
		executors: executors,
	}
}

// ScriptExecutionContext gives context to the execution of a plan script
type ScriptExecutionContext struct {
	ScriptName string
	Script     config.ShuttlePlanScript
	Project    config.ShuttleProjectContext
	Args       map[string]string
}

// ActionExecutionContext gives context to the execution of Actions in a script
type ActionExecutionContext struct {
	ScriptContext ScriptExecutionContext
	Action        config.ShuttleAction
	ActionIndex   int
}

// Execute is the command executor for the plan files
func (r *Registry) Execute(ctx context.Context, p config.ShuttleProjectContext, command string, args []string, validateArgs bool) error {
	script, ok := p.Scripts[command]
	if !ok {
		return errors.NewExitCode(2, "Script '%s' not found", command)
	}

	namedArgs, err := validateArguments(p, command, script.Args, args, validateArgs)
	if err != nil {
		return err
	}

	scriptContext := ScriptExecutionContext{
		ScriptName: command,
		Script:     script,
		Project:    p,
		Args:       namedArgs,
	}

	for actionIndex, action := range script.Actions {
		actionContext := ActionExecutionContext{
			ScriptContext: scriptContext,
			Action:        action,
			ActionIndex:   actionIndex,
		}
		err := r.executeAction(ctx, p.UI, actionContext)
		if err != nil {
			return err
		}
	}
	return nil
}

// validateArguments parses and validates args against available arguments in
// scriptArgs.
//
// All detectable constraints are checked before reporting to the UI.
func validateArguments(p config.ShuttleProjectContext, command string, scriptArgs []config.ShuttleScriptArgs, args []string, validateArgs bool) (map[string]string, error) {
	var validationErrors []validationError

	namedArgs, parsingErrors := validateArgFormat(args)
	validationErrors = append(validationErrors, parsingErrors...)
	if validateArgs {
		validationErrors = append(validationErrors, validateRequiredArgs(scriptArgs, namedArgs)...)
		validationErrors = append(validationErrors, validateUnknownArgs(scriptArgs, namedArgs)...)
	}
	if len(validationErrors) != 0 {
		sortValidationErrors(validationErrors)
		var s strings.Builder
		s.WriteString("Arguments not valid:\n")
		for _, e := range validationErrors {
			fmt.Fprintf(&s, " %s\n", e)
		}
		fmt.Fprintf(&s, "\n%s", expectedArgumentsHelp(command, scriptArgs))
		return nil, errors.NewExitCode(2, s.String())
	}
	return namedArgs, nil
}

type validationError struct {
	arg string
	err string
}

func (v validationError) String() string {
	return fmt.Sprintf("'%s' %s", v.arg, v.err)
}

func validateArgFormat(args []string) (map[string]string, []validationError) {
	var validationErrors []validationError
	namedArgs := map[string]string{}
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) < 2 {
			validationErrors = append(validationErrors, validationError{
				arg: arg,
				err: "not <argument>=<value>",
			})
			continue
		}
		namedArgs[parts[0]] = parts[1]
	}
	return namedArgs, validationErrors
}

func validateRequiredArgs(scriptArgs []config.ShuttleScriptArgs, args map[string]string) []validationError {
	var validationErrors []validationError
	for _, argSpec := range scriptArgs {
		if _, ok := args[argSpec.Name]; argSpec.Required && !ok {
			validationErrors = append(validationErrors, validationError{
				arg: argSpec.Name,
				err: "not supplied but is required",
			})
		}
	}
	return validationErrors
}

func validateUnknownArgs(scriptArgs []config.ShuttleScriptArgs, args map[string]string) []validationError {
	var validationErrors []validationError
	for namedArg := range args {
		found := false
		for _, arg := range scriptArgs {
			if arg.Name == namedArg {
				found = true
				break
			}
		}
		if !found {
			validationErrors = append(validationErrors, validationError{
				arg: namedArg,
				err: "unknown",
			})
		}
	}
	return validationErrors
}

func sortValidationErrors(errs []validationError) {
	sort.Slice(errs, func(i, j int) bool {
		return errs[i].arg < errs[j].arg
	})
}

func expectedArgumentsHelp(command string, args []config.ShuttleScriptArgs) string {
	var s strings.Builder
	fmt.Fprintf(&s, "Script '%s' accepts ", command)
	if len(args) == 0 {
		s.WriteString("no arguments.")
		return s.String()
	}
	s.WriteString("the following arguments:")
	for _, a := range args {
		fmt.Fprintf(&s, "\n  %s", a)
	}
	return s.String()
}

func (r *Registry) executeAction(ctx context.Context, ui *ui.UI, context ActionExecutionContext) error {
	for _, executor := range r.executors {
		handler, ok := executor(context.Action)
		if ok {
			return handler(ctx, ui, context)
		}
	}

	panic(fmt.Sprintf("Could not find an executor for %v.actions[%v]!", context.ScriptContext.ScriptName, context.ActionIndex))
}
