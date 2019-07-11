package executors

import (
	"fmt"
	"strings"

	"github.com/lunarway/shuttle/pkg/config"
)

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
func Execute(p config.ShuttleProjectContext, command string, args []string) {
	script, ok := p.Scripts[command]
	if !ok {
		p.UI.ExitWithErrorCode(2, "Script '%s' not found", command)
	}

	namedArgs := validateArguments(p, command, script.Args, args)

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
		executeAction(actionContext)
	}
}

// validateArguments parses and validates args against available arguments in
// scriptArgs.
//
// All detectable constraints are checked before reporting to the UI.
func validateArguments(p config.ShuttleProjectContext, command string, scriptArgs []config.ShuttleScriptArgs, args []string) map[string]string {
	var validationErrors []string

	namedArgs, parsingErrors := validateArgFormat(args)
	validationErrors = append(validationErrors, parsingErrors...)
	validationErrors = append(validationErrors, validateRequiredArgs(scriptArgs, namedArgs)...)
	validationErrors = append(validationErrors, validateUnknownArgs(scriptArgs, namedArgs)...)
	if len(validationErrors) != 0 {
		var s strings.Builder
		s.WriteString("Arguments not valid:\n")
		for _, e := range validationErrors {
			fmt.Fprintf(&s, " %s\n", e)
		}
		fmt.Fprintf(&s, "\n%s", expectedArgumentsHelp(command, scriptArgs))
		p.UI.ExitWithError(s.String())
	}
	return namedArgs
}

func validateArgFormat(args []string) (map[string]string, []string) {
	var validationErrors []string
	namedArgs := map[string]string{}
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) < 2 {
			validationErrors = append(validationErrors, fmt.Sprintf("'%s' not <argument>=<value>", arg))
			continue
		}
		namedArgs[parts[0]] = parts[1]
	}
	return namedArgs, validationErrors
}

func validateRequiredArgs(scriptArgs []config.ShuttleScriptArgs, args map[string]string) []string {
	var validationErrors []string
	for _, argSpec := range scriptArgs {
		if _, ok := args[argSpec.Name]; argSpec.Required && !ok {
			validationErrors = append(validationErrors, fmt.Sprintf("'%s' not supplied but is required", argSpec.Name))
		}
	}
	return validationErrors
}

func validateUnknownArgs(scriptArgs []config.ShuttleScriptArgs, args map[string]string) []string {
	var validationErrors []string
	for namedArg := range args {
		found := false
		for _, arg := range scriptArgs {
			if arg.Name == namedArg {
				found = true
				break
			}
		}
		if !found {
			validationErrors = append(validationErrors, fmt.Sprintf("'%s' unknown", namedArg))
		}
	}
	return validationErrors
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

func executeAction(context ActionExecutionContext) {
	for _, executor := range executors {
		if executor.match(context.Action) {
			executor.execute(context)
			return
		}
	}

	panic(fmt.Sprintf("Could not find an executor for %v.actions[%v]!", context.ScriptContext.ScriptName, context.ActionIndex))
}

type actionMatchFunc = func(config.ShuttleAction) bool
type actionExecutionFunc = func(ActionExecutionContext)

type actionExecutor struct {
	match   actionMatchFunc
	execute actionExecutionFunc
}

var executors = []actionExecutor{}

// AddExecutor taps a new executor into the script execution pipeline
func addExecutor(match actionMatchFunc, execute actionExecutionFunc) {
	executors = append(executors, actionExecutor{
		match:   match,
		execute: execute,
	})
}
