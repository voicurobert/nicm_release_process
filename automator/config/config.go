package config

import (
	"fmt"
	"github.com/bigkevmcd/go-configparser"
	"os"
	"path"
)

type parameters map[string]map[string]string

func getConfigPath() string {
	dir, _ := os.Getwd()
	return path.Join(dir, "nicm_paths.config")
}

func GetConfig() parameters {
	configFilePath := getConfigPath()
	config, err := configparser.NewConfigParserFromFile(configFilePath)
	if err != nil {
		panic(fmt.Sprintf("cannot read config file %s, error: %s", configFilePath, err.Error()))
	}

	configMap := make(parameters)

	for _, section := range config.Sections() {
		keyValue, _ := config.Items(section)
		configMap[section] = keyValue
	}

	return configMap
}

func GetCredentials() (string, string, string) {
	cfgMap := GetConfig()
	cfg, ok := cfgMap["server_config"]
	if ok {
		return cfg["ip"], cfg["username"], cfg["password"]
	}
	return "", "", ""
}
