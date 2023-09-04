// conf 配置加载器
package conf

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type ServerConf struct {
	Address string `yaml:"address"`
	Secret  string `yaml:"secret"`
	Post    string `yaml:"post"`
	Version string `yaml:"version"`
	Timeout string `yaml:"timeout"`
}

type Conf struct {
	Admin  []string   `yaml:"admin"`
	Server ServerConf `yaml:"server"`
}

var (
	config     Conf
	configLock sync.RWMutex
)

func Load(filePath string) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("无法加载配置文件：%w", err)
	}

	var newConfig Conf
	err = yaml.Unmarshal(bytes, &newConfig)
	if err != nil {
		return fmt.Errorf("无法解析配置文件：%w", err)
	}

	configLock.Lock()
	config = newConfig
	configLock.Unlock()

	return nil
}

func Get() Conf {
	configLock.RLock()
	defer configLock.RUnlock()

	return config
}

func Reload(filePath string) error {
	err := Load(filePath)
	if err != nil {
		return fmt.Errorf("无法重载配置文件：%w", err)
	}

	return nil
}
