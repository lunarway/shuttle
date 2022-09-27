package utility

import "github.com/lunarway/shuttle/pkg/config"

type ContextProvider func() (config.ShuttleProjectContext, error)
