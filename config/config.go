package config

import (
	"path"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	RawValue  interface{}
	Configers map[string]Config = make(map[string]Config)
)

type Config interface {
	Unmarshal(rawVal interface{}) error
}

func ParseConfig(configPath string, rawVal interface{}, watchChange bool) error {
	var dir, file string = path.Split(configPath)
	var fileName []string = strings.Split(file, ".")
	viper.SetConfigName(fileName[0])
	viper.SetConfigType(fileName[1])
	viper.AddConfigPath(dir)
	configer := Configers[fileName[1]]
	parseConfig(rawVal, configer)
	RawValue = rawVal

	if watchChange {
		go func() {
			viper.WatchConfig()
			viper.OnConfigChange(func(in fsnotify.Event) {
				parseConfig(rawVal, configer)
			})
		}()
	}

	return nil
}

func parseConfig(rawVal interface{}, configer Config) error {
	return configer.Unmarshal(rawVal)
}

// Register makes a configer available by the configer name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, configer Config) {
	if configer == nil {
		panic("config: Register configer is nil")
	}
	if _, ok := Configers[name]; ok {
		panic("config: Register called twice for configer " + name)
	}
	Configers[name] = configer
}
