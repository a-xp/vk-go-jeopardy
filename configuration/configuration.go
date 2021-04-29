package configuration

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Configuration struct {
	Mongo struct {
		Url     string
		Host    string
		Options string
		Timeout int
		Name    string
		User    string
		Pass    string
	}
	Http struct {
		ListenAddr string
		Mode       string
	}
	VkApp struct {
		Secret string
	}
	MockResponse bool `yaml:"mockResponse"`
}

func printConfigError(err error) {
	log.Printf("Failed to read configuration file %+v", err)
	os.Exit(2)
}

func createDefaultConfig() *Configuration {
	cfg := Configuration{}
	cfg.Http.ListenAddr = "0.0.0.0:9010"
	cfg.Http.Mode = "debug"
	cfg.MockResponse = true
	return &cfg
}

func LoadConfigFile() *Configuration {
	cfg := createDefaultConfig()
	f, err := os.Open("application.yml")
	if err != nil {
		printConfigError(err)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		printConfigError(err)
	}
	return cfg
}
