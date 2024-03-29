package externalservices

import (
	"encoding/json"
	"os"
	"path"
)

type Config struct {
	S3Config       S3Config       `json:"s3"`
	MastodonConfig MastodonConfig `json:"mastodon"`
}

type MastodonConfig struct {
	Server       string `json:"server"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	AccessToken  string `json:"accessToken"`
	Email        string `json:"email"`
	Password     string `json:"password"`
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

	configBytes, err := os.ReadFile(configPath)
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

func (m *MastodonConfig) Validate() bool {
	if m.Server == "" {
		return false
	}
	if m.ClientID == "" {
		return false
	}
	if m.ClientSecret == "" {
		return false
	}
	if m.AccessToken == "" {
		return false
	}
	if m.Email == "" {
		return false
	}
	if m.Password == "" {
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
