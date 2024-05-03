package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Http Http
}

type Http struct {
	Host string
	Port uint16
}

func Load(destination string) (Config, error) {
	content, err := os.ReadFile(destination)
	if err != nil {
		return Config{}, err
	}
	var result Config
	err = yaml.Unmarshal(content, &result)
	return result, err
}
