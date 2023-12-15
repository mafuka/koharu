package core

// Bot represents a QQ bot.
type Bot struct {
	*Config
	*Server
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

	RegisterLogger(b.Config.Log)

	RegisterEvent()
	RegisterMsgElem()

	if b.Config.isDefault {
		Log().Warn("Using the default configuration, some settings may not meet your expectations!")
	}
	if b.Config.Log.File == "console" {
		Log().Warn("Using console as log file, the logs will be lost when the console is closed!")
	}
	if b.Config.Log.Level == DebugLevel {
		Log().Warn("Using DEBUG log level! Tread with care, are you sure you dare?")
	}
	if len(b.Config.Admin) == 0 {
		Log().Warn("No bot admins assigned, console commands are disabled.")
	}
	if len(b.Config.Whitelist.Group) == 0 && len(b.Config.Whitelist.Friend) == 0 {
		Log().Warn("No whitelist set! Brace for impact, every message is now in play.")
	}

	b.Server = NewServer(b.Config.Server)

	b.Server.POST("/hook", Middlewares...)

	return b
}

// WithConfig is an option to set or override the configuration for the Bot.
func WithConfig(cfg *Config) Option {
	return func(b *Bot) {
		cfg.isDefault = false
		b.Config = cfg
	}
}

var Middlewares = []Middleware{Rx()}

func (b *Bot) Use(m Middleware) {
	Middlewares = append(Middlewares, m)
}

// Run starts the Bot.Server.
func (b *Bot) Run() error {
	Log().Info("Bot is listening for events on http://%s/hook.", b.Config.Server.Address)
	Log().Info("Bot will report behavior to http://%s.", b.Config.Server.Post)

	if err := b.Server.Run(); err != nil {
		return err
	}

	Log().Info("Good Dream.")
	return nil
}
