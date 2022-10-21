package config

import (
	"encoding/json"
	"os"
)

type GlesysHost struct {
	Domain     string   `json:"domain"`
	Subdomains []string `json:"subdomains"`
}
type GlesysCredentials struct {
	User string `json:"user"`
	Key  string `json:"key"`
}

type GlesysConfig struct {
	Hosts       []GlesysHost      `json:"hosts"`
	Credentials GlesysCredentials `json:"credentials"`
}

func Load(path string) (GlesysConfig, error) {
	cfg := GlesysConfig{}
	file, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&cfg)
	return cfg, err
}
