package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

const configFilename = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := fmt.Sprintf("%s%c%s", homedir, os.PathSeparator, configFilename)
	return configPath, nil
}

func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}
	c := Config{}
	if err := json.Unmarshal(configData, &c); err != nil {
		return Config{}, err
	}
	return c, nil
}

func (c *Config) write() error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	configData, err := json.Marshal(c)
	if err != nil {
		return err
	}
	os.WriteFile(configPath, configData, 0644)
	return nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUsername = username
	err := c.write()
	return err
}
