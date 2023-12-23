package core

import (
	"github.com/gin-contrib/pprof"
)

// Bot represents a QQ bot.
type Bot struct {
	*Config
	*Server

	Middlewares []Middleware

	withConfig bool
	withPProf  bool
}

// Option defines a function type for Bot options.
type Option func(*Bot)

// New initializes a new Bot instance with given options.
func New(options ...Option) *Bot {
	b := &Bot{
		Config:      DefaultConfig(),
		Server:      DefaultServer(),
		Middlewares: []Middleware{Rx()},
		withConfig:  false,
		withPProf:   false,
	}

	RegisterLogger(b.Config.Log)
	//RegisterEvent()
	//RegisterMsgElem()

	// apply any provided options
	for _, option := range options {
		option(b)
	}

	if !b.withConfig {
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

	return b
}

// WithConfig is an option to set or override the configuration for the Bot.
func WithConfig(cfg *Config) Option {
	return func(b *Bot) {
		b.withConfig = true
		b.Config = cfg
		b.Server = NewServer(cfg.Server)
	}
}

// WithPProf is a bot option to enable PProf.
func WithPProf() Option {
	return func(b *Bot) {
		b.withPProf = true
		pprof.Register(b.Engine)
	}
}

func (b *Bot) Use(m Middleware) {
	b.Middlewares = append(b.Middlewares, m)
}

// Run starts the Bot.Server.
func (b *Bot) Run() error {
	b.setupHook()

	Log().Info("Bot is listening for events on http://%s/hook.", b.httpServer.Addr)
	Log().Info("Bot will report behavior to %s.", b.Config.Server.Post)
	if b.withPProf {
		Log().Info("Bot has enabled PProf at http://%s/debug", b.httpServer.Addr)
	}

	if err := b.Server.Run(); err != nil {
		return err
	}

	Log().Info("Good Dream.")
	return nil
}

func (b *Bot) setupHook() {
	b.Server.POST("/hook", b.Middlewares...)
}
