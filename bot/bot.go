package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Bot represents a bot.
type Bot struct {
	Config *Config
	Server *Server
}

// Config holds all configurations for the Bot.
type Config struct {
	Admin  []int        `yaml:"admin"`
	Log    ZapConfig    `yaml:"log"`
	Server ServerConfig `yaml:"server"`
}

// DefaultConfig creates a new Config with default settings.
func DefaultConfig() *Config {
	return &Config{
		Admin: []int{},
		Log: ZapConfig{
			File:     "console",
			Level:    DebugLevel,
			MaxDays:  3,
			Compress: false,
		},
		Server: ServerConfig{
			Host:        "127.0.0.1",
			Port:        8080,
			Path:        "/onebot/v11/ws",
			AccessToken: "",
		},
	}
}

// Option defines a function type for Bot options.
type Option func(*Bot)

// WithLogger provides an option to set a custom logger.
func WithLogger(l Logger) Option {
	return func(b *Bot) {
		SetLogger(l)
	}
}

// New initializes a new Bot instance with given options.
func New(cfg *Config, options ...Option) *Bot {
	defaultLogger := NewZapLogger(cfg.Log)
	SetLogger(defaultLogger)

	b := &Bot{
		Config: cfg,
		Server: NewServer(cfg.Server),
	}

	// apply any provided options
	for _, option := range options {
		option(b)
	}

	return b
}

func (b *Bot) Run() error {
	// Bot start process begins
	Log().Debugf("Bot is starting up...")

	// Start the Server
	Log().Debugf("Starting the server...")
	if err := b.Server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	Log().Infof("Server is listening on ws://%s:%d%s",
		b.Config.Server.Host, b.Config.Server.Port, b.Config.Server.Path,
	)

	// Start other components...

	// Bot start process ends
	Log().Infof("Bot is up and running")

	// Wait for exit signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Bot stop process begins
	Log().Infof("Shutting down the bot...")

	// Stop the Server
	Log().Debugf("Shutting down the server...")
	if err := b.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}
	Log().Infof("Server gracefully closed")

	// Stop other components and clean up resources...

	// Bot stop process ends
	Log().Infof("Bot has been properly shut down")
	return nil
}
