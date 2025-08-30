package configLoader

import (
	"os"

	"gopkg.in/yaml.v3"
)

/**
 * Topic represents a Kafka topic configuration.
 */
type Topic struct {
	TopicName  string `yaml:"name"`
	GroupId    string `yaml:"groupId"`
	WorkerName string `yaml:"worker"`
}

/**
 * Config represents the configuration for Kafka consumer.
 */
type Config struct {
	Kafka struct {
		Host           string  `yaml:"host"`
		Port           int     `yaml:"port"`
		CommitInterval int     `yaml:"commitInterval"`
		Log            string  `yaml:"log"`
		Topics         []Topic `yaml:"topics"`
	} `yaml:"kafka"`
}

/**
 * LoadConfig loads the configuration from a YAML file.
 */
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
