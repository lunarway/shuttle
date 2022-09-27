package plugins

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
	"path"
	"sort"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/lunarway/shuttle/cmd/utility"
	"github.com/lunarway/shuttle/pkg/ui"
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"golang.org/x/mod/sumdb/dirhash"
)

func newLsCmd(uii *ui.UI, contextProvider utility.ContextProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use: "ls",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.ParseFlags(args); err != nil {
				return err
			}

			provider, err := contextProvider()
			if err != nil {
				return err
			}

			uii.Titleln("Plugins")
			packages := make([]string, 0)
			packagesMu := sync.Mutex{}
			wg := sync.WaitGroup{}
			wg.Add(len(provider.Config.Plugins))
			for k, v := range provider.Config.Plugins {
				go func(k string, v interface{}) {
					defer wg.Done()

					name, err := loadPlugin(cmd.Context(), k, v, provider.LocalShuttleDirectoryPath)
					if err != nil {
						panic(err)
					}

					packagesMu.Lock()
					defer packagesMu.Unlock()
					packages = append(packages, name)
				}(k, v)
			}
			wg.Wait()
			sort.Strings(packages)

			for _, p := range packages {
				uii.Infoln("  - %s", p)
			}

			return nil
		},
	}

	return cmd
}

func loadPlugin(ctx context.Context, pluginPath string, value interface{}, shuttleDir string) (string, error) {
	pluginDir := path.Join(shuttleDir, "plugins")
	if _, err := os.Stat(pluginDir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(pluginDir, 0777)
		if err != nil {
			return "", err
		}
	}

	fetchedPluginPath, pluginHash, err := fetchPlugin(ctx, pluginPath, pluginDir)
	if err != nil {
		return "", fmt.Errorf("could not fetch plugin: %w", err)
	}

	// TODO: Communicate with it. (Get its name)
	name, err := getName(ctx, fetchedPluginPath, pluginHash)
	if err != nil {
		return "", nil
	}

	return name, nil
}

func fetchPlugin(ctx context.Context, pluginPath string, pluginsDir string) (string, string, error) {
	p := path.Join(pluginsDir, "../../", pluginPath)

	if _, err := os.Stat(p); !os.IsNotExist(err) {
		p, pluginHash, err := copyPlugin(ctx, pluginPath, pluginsDir)
		if err != nil {
			return "", "", err
		}
		return p, pluginHash, nil
	} else if os.IsNotExist(err) {
		// TODO: the string is probably a url, fetch instead of aborting
		return "", "", fmt.Errorf("plugin path does not exist, aborting (%s)", pluginPath)
	}

	return "", "", errors.New("could not find plugin")
}

func copyPlugin(ctx context.Context, pluginsPath, pluginsDir string) (string, string, error) {
	hash, err := dirhash.HashDir(pluginsPath, "plugins", dirhash.Hash1)
	if err != nil {
		return "", "", err
	}
	hashB := sha256.Sum256([]byte(hash))
	hashS := hex.EncodeToString(hashB[:])

	p := path.Join(pluginsDir, hashS)

	return p, hashS, cp.Copy(pluginsPath, p)
}

func getName(ctx context.Context, pluginPath string, pluginHash string) (string, error) {
	logger := hclog.Default()
	logger.SetLevel(hclog.Error)

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "BASIC_PLUGIN",
			MagicCookieValue: "shuttle-plugin",
		},
		Plugins: map[string]plugin.Plugin{
			pluginHash: &ShuttlePlugin{},
		},
		Cmd:    exec.Command("sh", "-c", fmt.Sprintf("(cd %s; go run main.go %s)", pluginPath, pluginHash)),
		Logger: logger,
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		return "", err
	}

	raw, err := rpcClient.Dispense(pluginHash)
	if err != nil {
		return "", err
	}

	shuttlePlugin, ok := raw.(ShuttlePluginContract)
	if !ok {
		return "", errors.New("plugin is not of the right type")
	}

	name := shuttlePlugin.GetName()

	return name, nil
}

type ShuttlePluginContract interface {
	GetName() string
}

type ShuttleRPCPlugin struct {
	client *rpc.Client
}

func (sp *ShuttleRPCPlugin) GetName() string {
	var resp string
	err := sp.client.Call("Plugin.GetName", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

type ShuttleRPCServer struct {
	Impl ShuttlePluginContract
}

func (sp *ShuttleRPCServer) GetName(args interface{}, resp *string) error {
	*resp = sp.Impl.GetName()

	return nil
}

type ShuttlePlugin struct {
	Impl ShuttlePluginContract
}

// Client implements plugin.Plugin
func (*ShuttlePlugin) Client(m *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ShuttleRPCPlugin{
		client: c,
	}, nil
}

// Server implements plugin.Plugin
func (sp *ShuttlePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ShuttleRPCServer{
		Impl: sp.Impl,
	}, nil
}
