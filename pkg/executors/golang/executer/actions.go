package executer

type (
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

func (a *Actions) Execute(action string, fn func() error) (ran bool, err error) {
	if a == nil {
		return false, nil
	}

	if _, ok := a.Actions[action]; ok {
		return true, fn()
	}

	return true, nil
}
