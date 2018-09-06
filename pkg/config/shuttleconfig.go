package config

import (
	"io/ioutil"
	"log"
	"path"

	yaml "gopkg.in/yaml.v2"
)

type DynamicYaml = map[string]interface{}

// ShuttleConfig describes the actual config for each project
type ShuttleConfig struct {
	Plan      string      `yaml:"plan"`
	Variables DynamicYaml `yaml:"vars"`
}

// GetConf loads the ShuttleConfig from yaml file in the project path
func (c *ShuttleConfig) GetConf(projectPath string) *ShuttleConfig {
	var configPath = path.Join(projectPath, "shuttle.yaml")

	//log.Printf("configpath: %s", configPath)

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
