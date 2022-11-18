package configuration

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	LDAP LDAPConfig `yaml:"ldap"`
}

type LDAPConfig struct {
	URL     string   `yaml:"url"`
	BaseDN  string   `yaml:"base_dn,omitempty"`
	Filter  string   `yaml:"filter,omitempty"`
	BindDN  string   `yaml:"bind_dn,omitempty"`
	Login   string   `yaml:"login"`
	Secret  string   `yaml:"secret"`
	Filters []string `yaml:"filters,omitempty"`
	Exclude []string `yaml:"exclude,omitempty"`
}

// New initializes the configuration
func New(configFile string) *Config {
	c := new(Config)
	fileContent := readConfigFile(configFile)
	err := yaml.Unmarshal(*fileContent, c)
	if err != nil {
		log.Printf("unable to parse config file: %s", err)
	}

	return c
}

func readConfigFile(configFile string) *[]byte {
	var content []byte
	if _, err := os.Stat(configFile); err == nil {
		content, err = os.ReadFile(configFile)
		if err != nil {
			log.Fatalf("unable to read config file: %s", err)
		}
		return &content
	}
	return &content
}
