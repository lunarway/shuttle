package executors

import (
	"fmt"

	"bitbucket.org/LunarWay/shuttle/pkg/config"
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
	script := p.Plan.Scripts[command]

	namedArgs := map[string]string{}
	for i, argSpec := range script.Args {
		if len(args)-1 < i {
			panic(fmt.Sprintf("Required %v (position %v) is not supplied!", argSpec.Name, i))
		}
		namedArgs[argSpec.Name] = args[i]
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
