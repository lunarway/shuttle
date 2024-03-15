package extensions

import "path"

func getRemoteRegistryIndex() string {
	return "index"
}

func getRemoteRegistryExtensionPath(name string) string {
	return path.Join(getRemoteRegistryIndex(), name)
}

func getRemoteRegistryExtensionPathFile(name string) string {
	return path.Join(getRemoteRegistryExtensionPath(name), "shuttle-extension.json")
}
