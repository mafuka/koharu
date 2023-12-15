package core

// Config holds all configurations for the Bot.
type Config struct {
	isDefault bool
	Admin     []int           `yaml:"admin"`     // Administrator list
	Log       LoggerConfig    `yaml:"log"`       // Logger settings
	Server    ServerConfig    `yaml:"server"`    // Server settings
	Whitelist WhitelistConfig `yaml:"whitelist"` // Message whitelist
}

// DefaultConfig creates a new Config with default settings.
func DefaultConfig() *Config {
	return &Config{
		isDefault: true,
		Admin:     []int{},
		Log:       DefaultLoggerConfig(),
		Server:    DefaultServerConfig(),
		Whitelist: DefaultWhitelistConfig(),
	}
}

// ServerConfig holds server communication behavior.
type ServerConfig struct {
	Address string `yaml:"address"` // HTTP listener address
	Secret  string `yaml:"secret"`  // Authentication key
	Post    string `yaml:"post"`    // Reverse POST address
	Timeout int    `yaml:"timeout"` // Reverse HTTP timeout in seconds
}

// DefaultServerConfig provides a basic default ServerCfg.
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Address: "127.0.0.1:5701",
		Secret:  "",
		Post:    "http://127.0.0.1:5700",
		Timeout: 5,
	}
}

// WhitelistConfig defines whitelist chats for messaging.
type WhitelistConfig struct {
	Friend []int `yaml:"friend"` // Friend chats
	Group  []int `yaml:"group"`  // Group chats
}

// DefaultWhitelistConfig provides a basic default WhitelistCfg.
func DefaultWhitelistConfig() WhitelistConfig {
	return WhitelistConfig{
		Friend: []int{},
		Group:  []int{},
	}
}
