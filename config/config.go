package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Yaml ConfigYaml
}

type ConfigYaml struct {
	Host string `yaml:"host" env-default: "127.0.0.1"`
	Port string `yaml:"port" env-default: "8080"`
}

func (cfg *Config) ConfigPath() string {
	config_dir := cfg.ConfigDir()

	path := config_dir + "/" + "config.yml"
	return path
}

func (cfg *Config) ConfigDir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	path := homedir + "/" + ".config/" + "pmd"
	return path
}

func (cfg *Config) CreateConfig() (bool, error) {
	cfgPath := cfg.ConfigPath()

	config := ConfigYaml{
		Host: "127.0.0.1",
		Port: "8080",
	}

	// Create the config.
	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
		return false, err
	}

	os.MkdirAll(cfg.ConfigDir(), os.ModePerm)
	err = ioutil.WriteFile(cfgPath, yamlData, 0644)
	if err != nil {
		return false, err
	}

	return true, nil
}
