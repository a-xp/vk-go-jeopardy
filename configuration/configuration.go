package configuration

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Configuration struct {
	Mongo struct {
		Scheme  string
		Host    string
		Options string
		Timeout int
		Name    string
		User    string
		Pass    string
	}
	Http struct {
		ListenAddr string `yaml:"listenAddr"`
		Mode       string `yaml:"mode"`
		PublicAddr string `yaml:"publicAddr"`
	}
	VkApp struct {
		Secret string
		Key    string
		Url    string
	}
	RootPass        string `yaml:"rootPass"`
	MockResponse    bool   `yaml:"mockResponse"`
	ValidateRequest bool   `yaml:"validateRequest"`
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
	cfg.ValidateRequest = true
	cfg.Mongo.Scheme = "mongo+srv"
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
	decoder.SetStrict(true)
	err = decoder.Decode(cfg)
	if err != nil {
		printConfigError(err)
	}
	return cfg
}
