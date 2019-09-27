package config

type ShuttleActionStep struct {
	Name string            `yaml:"name"`
	Args map[string]string `yaml:"args"`
}
