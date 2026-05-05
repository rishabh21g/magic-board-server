package config

import "github.com/rishabh21g/magic-board/pkg/utils"

type Config struct {
	Port          string
	RedisAddress  string
	RedisPassword string
	RedisUser     string
}

func LoadConfig() *Config {
	return &Config{
		Port:          utils.GetEnv("PORT", "8080"),
		RedisAddress:  utils.GetEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: utils.GetEnv("REDIS_PASSWORD", ""),
		RedisUser:     utils.GetEnv("REDIS_USER", "default"),
	}
}
