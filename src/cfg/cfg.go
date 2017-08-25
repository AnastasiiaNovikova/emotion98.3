package cfg

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func GetConfigDir() string {
	return "../config"
}

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

type App struct {
	Cognitron struct {
		MaxJobs string `yaml:"max_jobs"`
	}
}

var app App

func GetApp() App {
	return app
}

func init() {
	if err := GetYamlConfig("app", &app); err != nil {
		log.Fatalf("can't read app config: %s", err)
	}
}
