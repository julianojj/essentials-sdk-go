package config

import _ "github.com/joho/godotenv/autoload"

type Config struct{}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) GetConfig() *Config {
	return c
}
