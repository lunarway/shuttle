package ui

import (
	"fmt"
	"os"
)

// UI is the abstraction of handling terminal output for shuttle
type UI struct {
	EffectiveLevel Level
	DefaultLevel   Level
	UserLevel      Level
	UserLevelSet   bool
}

// Create doc
func Create() UI {
	return UI{
		EffectiveLevel: LevelInfo,
		DefaultLevel:   LevelInfo,
		UserLevelSet:   false,
	}
}

// SetUserLevel doc
func (ui *UI) SetUserLevel(level Level) UI {
	return UI{
		EffectiveLevel: level,
		DefaultLevel:   ui.DefaultLevel,
		UserLevel:      level,
		UserLevelSet:   true,
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
	}
}

// VerboseLn doc
func (ui *UI) VerboseLn(msg string, args ...interface{}) {
	if ui.EffectiveLevel.OutputIsIncluded(LevelVerbose) {
		fmt.Println(fmt.Sprintf(msg, args...))
	}
}

// InfoLn doc
func (ui *UI) InfoLn(msg string, args ...interface{}) {
	if ui.EffectiveLevel.OutputIsIncluded(LevelInfo) {
		fmt.Println(fmt.Sprintf(msg, args...))
	}
}

// EmphasizeInfoLn doc
func (ui *UI) EmphasizeInfoLn(msg string, args ...interface{}) {
	if ui.EffectiveLevel.OutputIsIncluded(LevelInfo) {
		ui.InfoLn("\x1b[032;1m"+msg+"\x1b[0m", args...)
	}
}

// TitleLn doc
func (ui *UI) TitleLn(msg string, args ...interface{}) {
	ui.InfoLn("\x1b[1m"+msg+"\x1b[0m", args...)
}

// ErrorLn doc
func (ui *UI) ErrorLn(msg string, args ...interface{}) {
	if ui.EffectiveLevel.OutputIsIncluded(LevelInfo) {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("\x1b[31;1m"+msg+"\x1b[0m", args...))
	}
}

// ExitWithError doc
func (ui *UI) ExitWithError(msg string, args ...interface{}) {
	ui.ExitWithErrorCode(1, msg, args...)
}

// ExitWithErrorCode doc
func (ui *UI) ExitWithErrorCode(code int, msg string, args ...interface{}) {
	ui.ErrorLn("shuttle failed\n"+msg, args...)
	os.Exit(code)
}

// CheckIfError doc
func (ui *UI) CheckIfError(err error) {
	if err == nil {
		return
	}
	ui.ExitWithError(fmt.Sprintf("error: %s", err))
}
