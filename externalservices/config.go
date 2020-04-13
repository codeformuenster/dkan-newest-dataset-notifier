package externalservices

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	TwitterConfig TwitterConfig `json:"twitter"`
}

type TwitterConfig struct {
	ConsumerKey    string `json:"consumerKey"`
	ConsumerSecret string `json:"consumerSecret"`
	AccessToken    string `json:"accessToken"`
	AccessSecret   string `json:"accessSecret"`
}

func FromFile(configPath string) (*Config, error) {
	if configPath == "" {
		dir, err := os.Getwd()
		if err != nil {
			return &Config{}, nil
		}

		configPath = path.Join(dir, "config.json")
	}

	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return &Config{}, err
	}
	c := Config{}
	err = json.Unmarshal(configBytes, &c)
	if err != nil {
		return &Config{}, err
	}

	return &c, nil
}

func (t *TwitterConfig) Validate() bool {
	if t.ConsumerKey == "" {
		return false
	}
	if t.ConsumerSecret == "" {
		return false
	}
	if t.AccessToken == "" {
		return false
	}
	if t.AccessSecret == "" {
		return false
	}
	return true
}
