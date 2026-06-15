package config

import (
	"fmt"
	"os"
)

type Config struct {
	BotToken  string
	Channel   string
	ProxyAddr string
	KeysFile  string
	Port      string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		BotToken:  os.Getenv("BOT_TOKEN"),
		Channel:   os.Getenv("CHANNEL"),
		ProxyAddr: os.Getenv("PROXY_ADDR"),
		KeysFile:  os.Getenv("KEYS_FILE"),
		Port:      os.Getenv("PORT"),
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.BotToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN is required")
	}
	if cfg.Channel == "" {
		return nil, fmt.Errorf("CHANNEL is required")
	}
	if cfg.KeysFile == "" {
		return nil, fmt.Errorf("KEYS_FILE is required")
	}
	return cfg, nil
}
