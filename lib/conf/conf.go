// conf 配置加载器
package conf

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Conf struct {
	Admin     []int64    `yaml:"admin"`
	Server    ServerConf `yaml:"server"`
	Whitelist Whitelist  `yaml:"whitelist"`
}

type ServerConf struct {
	Address string `yaml:"address"`
	Secret  string `yaml:"secret"`
	Post    string `yaml:"post"`
	Version string `yaml:"version"`
	Timeout string `yaml:"timeout"`
}

type Whitelist struct {
	Private []int64 `yaml:"private"`
	Group   []int64 `yaml:"group"`
}

var (
	conf     Conf
	confLock sync.RWMutex
)

func Load(filePath string) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read configuration: %w", err)
	}

	var newConf Conf
	err = yaml.Unmarshal(bytes, &newConf)
	if err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

	confLock.Lock()
	conf = newConf
	confLock.Unlock()

	return nil
}

func Get() Conf {
	confLock.RLock()
	defer confLock.RUnlock()

	return conf
}

func Reload(filePath string) error {
	err := Load(filePath)
	if err != nil {
		return fmt.Errorf("failed to load configuration：%w", err)
	}

	return nil
}
