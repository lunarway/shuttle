package main

import (
	"os"

	"github.com/hashicorp/go-plugin"
	"github.com/lunarway/shuttle/cmd/plugins"
)

type ShuttlePlugin struct{}

var _ plugins.ShuttlePluginContract = &ShuttlePlugin{}

func (sp *ShuttlePlugin) GetName() string {
	return "plugin-two"
}

func main() {
	args := os.Args
	if len(args) != 2 {
		panic("args not supplied correctly")
	}

	sp := &ShuttlePlugin{}

	pluginMap := map[string]plugin.Plugin{
		args[1]: &plugins.ShuttlePlugin{Impl: sp},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "BASIC_PLUGIN",
			MagicCookieValue: "shuttle-plugin",
		},
		Plugins: pluginMap,
	})
}
