package plan

import (
	"github.com/emilingerslev/shuttle/pkg/config"
	"github.com/emilingerslev/shuttle/pkg/docker"
)

// Execute is the command executor for the plan files
func Execute(p config.ShuttlePlan, command string) string {
	script := p.Configuration.Scripts[command]

	switch {
	case script.Shell != "":
		return executeShell(p, command, script)
	case script.Dockerfile != "":
		return executeDocker(p, command, script)
	}
	panic("No valid execution found for command '" + command + "'!")
}

func executeDocker(p config.ShuttlePlan, command string, s config.ShuttlePlanScript) string {
	docker.Build(p, command, s)
	return command + "> Docker: Executing docker - " + s.Dockerfile
}
