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
		p.UI.ExitWithErrorCode(2, "No script found called '%s'", command)
	}
	namedArgs := map[string]string{}
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) < 2 {
			p.UI.ExitWithError("Could not parse `shuttle run %s %s`, because '%s' was expected to be on the `<option>=<value>` form, but wasn't!.\nScript '%s' expects arguments:\n%v", command, strings.Join(args, " "), arg, command, script.Args)
		}
		namedArgs[parts[0]] = parts[1]
	}

	for _, argSpec := range script.Args {
		if _, ok := namedArgs[argSpec.Name]; argSpec.Required && !ok {
			p.UI.ExitWithError("Required argument `%s` was not supplied!", argSpec.Name) // TODO: Add expected arguments
		}
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
		executeAction(actionContext)
	}
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
