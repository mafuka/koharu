package core

import (
	"bufio"
	"golang.org/x/xerrors"
	"os"
)

func init() {
	NewLogger(DefaultConfig().Log)
}

// Bot represents a bot.
type Bot struct {
	*Config

	Admin []int
	*Server
	*Client
	Middlewares []Middleware
}

// Config holds all configurations for the Bot.
type Config struct {
	isDefault bool

	Log    LogConfig    `yaml:"log"`
	Server ServerConfig `yaml:"server"`
	Client ClientConfig `yaml:"client"`

	Admin []int `yaml:"admin"`
}

// DefaultConfig creates a new Config with default settings.
func DefaultConfig() *Config {
	return &Config{
		isDefault: true,
		Log: LogConfig{
			File:     "console",
			Level:    InfoLevel,
			MaxDays:  3,
			Compress: false,
			JSON:     false,
		},
		Server: ServerConfig{
			Address: "127.0.0.1:5701",
			Secret:  "",
		},
		Client: ClientConfig{
			ID:      0,
			Key:     "nil",
			Address: "http://127.0.0.1:5700",
			Timeout: 5,
		},
		Admin: []int{},
	}
}

// Option defines a function type for Bot options.
type Option func(*Bot)

// New initializes a new Bot instance with given options.
func New(options ...Option) *Bot {
	b := &Bot{}
	b.Config = DefaultConfig()
	b.Server = NewServer(b.Config.Server)
	b.Client = NewClient(b.Config.Client)

	// apply any provided options
	for _, option := range options {
		option(b)
	}
	Log().Debug("All options have been applied")

	// Validate configuration and log warnings.
	validateConfig(b.Config)

	return b
}

// validateConfig checks the bot's configuration and logs warnings.
func validateConfig(cfg *Config) {
	if cfg.isDefault {
		Log().Error("No configuration explicitly specified, the bot will work abnormally!")
	}
	if cfg.Log.File == "console" {
		Log().Warn("Using console as log file!")
	}
	if cfg.Log.Level == DebugLevel {
		Log().Warn("Using DEBUG log level!")
	}
	if len(cfg.Admin) == 0 {
		Log().Warn("No bot admins assigned!")
	}
}

// WithConfig is an Option to set or override the configuration for the Bot.
func WithConfig(cfg *Config) Option {
	return func(b *Bot) {
		cfg.isDefault = false
		b.Config = cfg
		b.Server = NewServer(cfg.Server)
		b.Client = NewClient(cfg.Client)
		Log().Update(cfg.Log)
		Log().Debug("Option WithConfig has been applied")
	}
}

// WithPProf is a Bot Option to enable PProf.
func WithPProf() Option {
	return func(b *Bot) {
		b.PProf()
		Log().Info("PProf enabled at http://%s/debug", b.httpServer.Addr)
		Log().Debug("Option WithPProf has been applied")
	}
}

// Use appends a Middleware to the bot.
func (b *Bot) Use(m Middleware) {
	b.Middlewares = append(b.Middlewares, m)
}

// Run starts the Bot.Server.
func (b *Bot) Run() error {
	if err := startClient(b); err != nil {
		Log().Error("Client startup failed, bot will only be able to receive events:\n%+v", err)
	}

	startHook(b)

	Log().Info("Bot is up and running")
	if err := b.Server.Run(); err != nil {
		Log().Error("Failed to start server: %v", err)
		return err
	}

	Log().Info("Bot shutdown gracefully")

	cleanup(b)
	return nil
}

func startClient(b *Bot) error {
	Log().Info("Is Mirai up and logged in? (Press any key to continue...)")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	Log().Info("Start preparing the client")

	c := b.Client
	Log().Debug("Getting a session for the client %d", c.ID)

	Log().Debug("Validating client %d", c.ID)
	session, err := c.Verify()
	if err != nil {
		return xerrors.Errorf("Unable to verify client: %w", err)
	}

	Log().Debug("Binding session %s to client %d", session, c.ID)
	if err := c.Bind(session); err != nil {
		return xerrors.Errorf("Unable to bind session to client: %w", err)
	}

	Log().Info("Client is logged in as %d", c.ID)
	Log().Info("Server will report behavior to %s", c.Address)
	return nil
}

func startHook(b *Bot) {
	b.Server.POST("/hook", b.Middlewares...)
	Log().Info("Bot is listening for events on http://%s/hook", b.httpServer.Addr)
}

func cleanup(b *Bot) {
	Log().Info("Starts cleanup operations")
	c := b.Client
	if c.session != "" {
		Log().Debug("Releasing session %s from client %d", c.session, c.ID)
		if err := c.Release(c.session); err != nil {
			Log().Error("Unable to release client session:\n%+v", err)
			return
		}
		Log().Info("Client session released.")
	}
	Log().Info("Cleanup complete")
}
