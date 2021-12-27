package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

const (
	configPath = "config.yaml"
)

type Config struct {
	Username         string `yaml:"username"`
	LobbyName        string `yaml:"lobbyName"`
	WalletPassphrase string `yaml:"walletPassphrase"`
	WalletPassword   string `yaml:"walletPassword"`
	AppUrl           string `yaml:"appUrl"`
}

var Cfg *Config

func init() {
	config := &Config{}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Fatal(err)
	}

	Cfg = config
	fmt.Println("config loaded")
}
