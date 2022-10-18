package config

import (
	"encoding/json"
	"os"
)

func LoadHosts(path string) (map[string][]string, error) {
	hosts := map[string][]string{}
	file, err := os.Open(path)
	if err != nil {
		return hosts, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&hosts)
	return hosts, err
}
