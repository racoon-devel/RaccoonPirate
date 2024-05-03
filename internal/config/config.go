package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Http      Http
	Discovery Discovery
}

type Http struct {
	Host string
	Port uint16
}

type Discovery struct {
	Identity string
	Scheme   string
	Host     string
	Port     uint16
	Path     string
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
