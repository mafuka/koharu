package bot

import (
	"github.com/lesismal/nbio/nbhttp"
)

// Bot represents a bot.
type Bot struct {
	*Logger
	Engine *nbhttp.Engine
}

// Config holds all configurations for the Bot.
type Config struct {
	Admin        []int `yaml:"admin"`
	LogConfig    `yaml:"log"`
	ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Host        string `yaml:"host"`
	Port        int32  `yaml:"port"`
	Path        string `yaml:"path"`
	AccessToken string `yaml:"access_token"`
}

// DefaultConfig creates a new Config with default settings.
func DefaultConfig() *Config {
	return &Config{
		Admin: []int{},
		LogConfig: LogConfig{
			File:     "console",
			Level:    DebugLevel,
			MaxDays:  3,
			Compress: false,
		},
		ServerConfig: ServerConfig{
			Host:        "127.0.0.1",
			Port:        8080,
			Path:        "/onebot/v11/ws",
			AccessToken: "",
		},
	}
}

// Option defines a function type for Bot options.
type Option func(*Bot)

// New initializes a new Bot instance with given options.
func New(cfg Config, options ...Option) *Bot {
	b := &Bot{}
	b.Logger = NewLogger(cfg.LogConfig)
	b.Engine = nbhttp.NewEngine(nbhttp.Config{})

	// apply any provided options
	for _, option := range options {
		option(b)
	}

	return b
}
