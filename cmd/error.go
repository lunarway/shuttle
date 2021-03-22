package cmd

import (
	"errors"
	"os"

	shuttleerrors "github.com/lunarway/shuttle/pkg/errors"
)

func checkError(err error) {
	if err == nil {
		return
	}
	var exitCode *shuttleerrors.ExitCode
	if errors.As(err, &exitCode) {
		uii.Errorln("shuttle failed\n%s", exitCode.Message)
		os.Exit(exitCode.Code)
	}
	uii.Errorln("shuttle failed\nError: %s", err)
	os.Exit(1)
}
