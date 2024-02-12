package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	ConfigurationFilePathFrom string = "./src/configuration/"
	ConfigurationFileName     string = "config.json"
)

type Configuration struct {
	IsDevEnvironment bool `json:"isDevEnvironment"`
}

func LoadConfiguration() (Configuration, error) {
	config := getConfigurationWithDefaults()

	configFile, err := os.Open(fmt.Sprintf("%s%s", ConfigurationFilePathFrom, ConfigurationFileName))

	if err != nil {
		return config, err
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParseErr := jsonParser.Decode(&config)

	configFile.Close()
	return config, jsonParseErr
}

// Creates a Configuration struct and sets fields to their default values
func getConfigurationWithDefaults() Configuration {
	return Configuration{
		IsDevEnvironment: true,
	}
}
