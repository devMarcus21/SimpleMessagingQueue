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

type LoggingConfiguration struct {
	LoggerType string `json:"loggerType"` // TODO make this a strong type somehow
}

type BatchingConfiguration struct {
	MaxBatchPushSize int `json:"maxBatchPushSize"`
	MaxBatchReadSize int `json:"maxBatchReadSize"`
}

type Configuration struct {
	Port             int                   `json:"port"`
	IsDevEnvironment bool                  `json:"isDevEnvironment"`
	Logging          LoggingConfiguration  `json:"logging"` // Can be "text" or "json"
	Batching         BatchingConfiguration `json:"batching"`
}

func LoadConfiguration() (Configuration, error) {
	config := NewConfigurationWithDefaults()

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
func NewConfigurationWithDefaults() Configuration {
	return Configuration{
		Port:             8080,
		IsDevEnvironment: true,
		Logging: LoggingConfiguration{
			LoggerType: "text",
		},
		Batching: BatchingConfiguration{
			MaxBatchPushSize: 10,
			MaxBatchReadSize: 10,
		},
	}
}
