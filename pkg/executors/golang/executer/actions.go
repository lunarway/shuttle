package executer

type (
	// Actions represents all the possible commands to be sent to the golang actions binaries.
	// Such as `shuttle run daggerbuild --arg something`, in this daggerbuild would be the name, and arg being an Arg for said action
	Actions struct {
		Actions map[string]Action `json:"actions"`
	}

	Action struct {
		Args []ActionArg `json:"args"`
	}

	ActionArg struct {
		Name string `json:"name"`
	}
)

func NewActions() *Actions {
	return &Actions{
		Actions: make(map[string]Action, 0),
	}
}

// Merge exists to combine multiple Actions from a variety of binaries into one.
// This allows a single set of actions to represent all possible actions by shuttle in a given context
func (a *Actions) Merge(other ...*Actions) *Actions {
	for _, actions := range other {
		if actions == nil {
			continue
		}
		for name, action := range actions.Actions {
			a.Actions[name] = action
		}
	}

	return a
}

// Execute can execute a single action given a name, and a closure to handle any actual execution
func (a *Actions) Execute(action string, fn func() error) (ran bool, err error) {
	if a == nil {
		return false, nil
	}

	if _, ok := a.Actions[action]; ok {
		return true, fn()
	}

	return false, nil
}
