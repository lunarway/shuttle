package executors

import (
	"fmt"
	"github.com/lunarway/shuttle/pkg/config"
	"time"
)

// executeStep runs another shuttle step
func executeStep(context ActionExecutionContext) {
	//shuttlePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	start := time.Now()
	fmt.Println("--- Run " + context.Action.Step.Name)

	var args []string
	for argTo, argFrom := range context.Action.Step.Args  {
		for k,v := range context.ScriptContext.Args {
			if argFrom == k {
				args = append(args, fmt.Sprintf("%s=%s", argTo,v))
			}
		}
	}

	Execute(context.ScriptContext.Project, context.Action.Step.Name, args)
	fmt.Println(fmt.Sprintf("completed in %vs", time.Now().Sub(start).Seconds()))
	fmt.Println("")
}

func init() {
	addExecutor(func(action config.ShuttleAction) bool {
		return action.Step.Name != ""
	}, executeStep)
}
