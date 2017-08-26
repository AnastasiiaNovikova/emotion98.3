package cfg

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// GetConfigDir ectory
func GetConfigDir() string {
	return "../config"
}

// GetEnv inronment (development as default)
func GetEnv() string {
	env := os.Getenv("GO_ENV")
	if env != "" {
		return env
	}

	if flag.Lookup("test.v") != nil {
		return "test"
	}

	return "development"
}

// GetYamlConfig from file to some object
func GetYamlConfig(name string, config interface{}) error {
	configPath := fmt.Sprintf("%s/%s.yaml", GetConfigDir(), name)
	configContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("can't read config %q: %s", configPath, err)
	}

	if err = yaml.Unmarshal(configContent, config); err != nil {
		return fmt.Errorf("invalid yaml in config %q: %s", configPath, err)
	}

	return nil
}

// App config
type App struct {
	Cognitron struct {
		MaxJobs int `yaml:"max_jobs"`
		Timeout int `yaml:"timeout"`
	}
}

var app App

// GetApp configuration
func GetApp() App {
	return app
}

func init() {
	if err := GetYamlConfig("app", &app); err != nil {
		log.Fatalf("can't read app config: %s", err)
	}
}
