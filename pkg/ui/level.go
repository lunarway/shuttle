package ui

// Level specifies the level of output that commands should print
type Level string

// LevelVerbose includes all output
const (
	LevelVerbose Level = "Verbose"
	LevelInfo    Level = "Info"
	LevelError   Level = "Error"
	LevelSilent  Level = "Silent"
)

var levelMap = map[Level]int{
	LevelVerbose: 3,
	LevelInfo:    2,
	LevelError:   1,
	LevelSilent:  0,
}

// OutputIsIncluded returns true if levelA specifies that levelB should be included for output
func (levelA Level) OutputIsIncluded(levelB Level) bool {
	return levelMap[levelA] >= levelMap[levelB]
}
