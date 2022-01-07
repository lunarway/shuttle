package ui

import (
	"fmt"
	"io"
	"os"
)

// UI is the abstraction of handling terminal output for shuttle
type UI struct {
	EffectiveLevel Level
	DefaultLevel   Level
	UserLevel      Level
	UserLevelSet   bool
	Out            io.Writer
	Err            io.Writer
}

// Create doc
func Create() UI {
	return UI{
		EffectiveLevel: LevelInfo,
		DefaultLevel:   LevelInfo,
		UserLevelSet:   false,
		Out:            os.Stdout,
		Err:            os.Stderr,
	}
}

func (ui *UI) SetOutput(out io.Writer, err io.Writer) {
	ui.Out = out
	ui.Err = err
}

// SetUserLevel doc
func (ui *UI) SetUserLevel(level Level) UI {
	return UI{
		EffectiveLevel: level,
		DefaultLevel:   ui.DefaultLevel,
		UserLevel:      level,
		UserLevelSet:   true,
		Out:            ui.Out,
		Err:            ui.Err,
	}
}

// SetContext doc
func (ui *UI) SetContext(level Level) UI {
	var effectiveLevel Level
	if ui.UserLevelSet {
		effectiveLevel = ui.UserLevel
	} else {
		effectiveLevel = level
	}

	return UI{
		EffectiveLevel: effectiveLevel,
		DefaultLevel:   level,
		UserLevel:      ui.UserLevel,
		UserLevelSet:   ui.UserLevelSet,
		Out:            ui.Out,
		Err:            ui.Err,
	}
}

// Verboseln prints a formatted verbose message line.
func (ui *UI) Verboseln(format string, args ...interface{}) {
	if ui.EffectiveLevel.OutputIsIncluded(LevelVerbose) {
		fmt.Fprintln(ui.Out, fmt.Sprintf(format, args...))
	}
}

// Infoln prints a formatted info message line.
func (ui *UI) Infoln(format string, args ...interface{}) {
	if ui.EffectiveLevel.OutputIsIncluded(LevelInfo) {
		fmt.Fprintln(ui.Out, fmt.Sprintf(format, args...))
	}
}

func (ui *UI) EmphasizeInfoln(format string, args ...interface{}) {
	if ui.EffectiveLevel.OutputIsIncluded(LevelInfo) {
		fmt.Fprintf(ui.Out, "\x1b[032;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
	}
}

// Titleln doc
func (ui *UI) Titleln(format string, args ...interface{}) {
	ui.Infoln("\x1b[1m%s\x1b[0m", fmt.Sprintf(format, args...))
}

// Errorln doc
func (ui *UI) Errorln(format string, args ...interface{}) {
	if ui.EffectiveLevel.OutputIsIncluded(LevelError) {
		fmt.Fprintf(ui.Err, "\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
	}
}
