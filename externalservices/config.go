package externalservices

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	TwitterConfig TwitterConfig `json:"twitter"`
	S3Config      S3Config      `json:"s3"`
}

type TwitterConfig struct {
	ConsumerKey    string `json:"consumerKey"`
	ConsumerSecret string `json:"consumerSecret"`
	AccessToken    string `json:"accessToken"`
	AccessSecret   string `json:"accessSecret"`
}

type S3Config struct {
	Region          string `json:"region"`
	Endpoint        string `json:"endpoint"`
	Bucket          string `json:"bucket"`
	Path            string `json:"path"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
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

func (s *S3Config) Validate() bool {
	if s.Region == "" {
		return false
	}
	if s.Endpoint == "" {
		return false
	}
	if s.Bucket == "" {
		return false
	}
	if s.Path == "" {
		return false
	}
	if s.AccessKeyID == "" {
		return false
	}
	if s.SecretAccessKey == "" {
		return false
	}

	return true
}
