package configs

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`
	WeatherAPI struct {
		Key string `yaml:"key"`
	} `yaml:"weather_api"`
	DB struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
		Name string `yaml:"name"`
	} `yaml:"db"`
	SMTP struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		From     string `yaml:"from"`
	} `yaml:"smtp_config"`
	Token struct {
		Secret string `yaml:"secret"`
	} `yaml:"token"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	expanded := os.ExpandEnv(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
