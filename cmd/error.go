package cmd

import (
	stdcontext "context"
	"errors"
	"os"

	shuttleerrors "github.com/lunarway/shuttle/pkg/errors"
	"github.com/lunarway/shuttle/pkg/ui"
)

func checkError(uii *ui.UI, err error) {
	if err == nil {
		return
	}
	var exitCode *shuttleerrors.ExitCode
	if errors.As(err, &exitCode) {
		uii.Errorln("shuttle failed\n%s", exitCode.Message)
		os.Exit(exitCode.Code)
	}
	if errors.Is(err, stdcontext.Canceled) {
		uii.Errorln("Operation cancelled")
		os.Exit(2)
	}
	uii.Errorln("shuttle failed\nError: %s", err)
	os.Exit(1)
}
