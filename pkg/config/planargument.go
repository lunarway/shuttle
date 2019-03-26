package config

import (
	"github.com/lunarway/shuttle/pkg/git"
	"os"
	"path"
	"strings"
)

func isPlanArgumentAFilePlan(planArgument string) bool {
	return strings.HasPrefix(planArgument, "/") || strings.HasPrefix(planArgument, "./") || strings.HasPrefix(planArgument, "../")
}

func isPlanArgumentAPlan(planArgument string) bool {
	return planArgument != "" && (git.IsGitPlan(planArgument) || isPlanArgumentAFilePlan(planArgument))
}

func getPlanFromPlanArgument(planArgument string) string {
	switch {
	case isPlanArgumentAFilePlan(planArgument) && isFilePath(planArgument, true):
		return planArgument
	case isPlanArgumentAFilePlan(planArgument) && isFilePath(planArgument, false):
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		return path.Join(wd, planArgument)
	default:
		return planArgument
	}
}
