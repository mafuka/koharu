package bot

// init creates a default Logger when calling the package.
func init() { NewLogger(DefaultConfig().LogConfig) }

// Bot represents a bot.
type Bot struct {
	*Config

	// httpClient *http.Client
}

// Config holds all configurations for the Bot.
type Config struct {
	isDefault bool

	Admin      []int `yaml:"admin"`
	LogConfig  `yaml:"log"`
	AuthConfig `yaml:"upstream"`
}

// DefaultConfig creates a new Config with default settings.
func DefaultConfig() *Config {
	return &Config{
		isDefault: true,
		Admin:     []int{},
		LogConfig: LogConfig{
			File:     "console",
			Level:    DebugLevel,
			MaxDays:  3,
			Compress: false,
		},
		AuthConfig: AuthConfig{
			APPID:     0,
			APPSecret: "",
		},
	}
}

// Option defines a function type for Bot options.
type Option func(*Bot)

// New initializes a new Bot instance with given options.
func New(options ...Option) *Bot {
	b := &Bot{}
	b.Config = DefaultConfig()

	// apply any provided options
	for _, option := range options {
		option(b)
	}

	// Validate configuration and log warnings.
	// validateConfig(b.Config)

	return b
}

// WithConfig is an Option to set or override the configuration for the Bot.
func WithConfig(cfg *Config) Option {
	return func(b *Bot) {
		cfg.isDefault = false
		b.Config = cfg
		Log().ReInit(cfg.LogConfig)
		Log().Debug("Option WithConfig has been applied")
	}
}
