package config

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"google.golang.org/appengine/log"
)

var (
	Configer IConfiger
)

type IConfiger interface {
	unmarshal() error
	loadData() error
	parse() error
}

type BaseConfiger struct {
	WatchChange bool
	RawVal      interface{}
}

func (bConfiger *BaseConfiger) unmarshal() error {
	return viper.Unmarshal(bConfiger.RawVal)
}

type FileConfiger struct {
	BaseConfiger
	ConfigName string
	ConfigType string
	ConfigPath string
}

func (fConfiger *FileConfiger) loadData() error {
	viper.SetConfigName(fConfiger.ConfigName)
	viper.SetConfigType(fConfiger.ConfigType)
	viper.AddConfigPath(fConfiger.ConfigPath)
	return viper.ReadInConfig()
}

func (fConfiger *FileConfiger) parse() error {
	var err = fConfiger.unmarshal()
	if err != nil {
		return err
	}

	if fConfiger.WatchChange {
		go func(configer IConfiger) {
			viper.WatchConfig()
			viper.OnConfigChange(func(in fsnotify.Event) {
				_ = configer.unmarshal()
			})
		}(fConfiger)
	}

	return nil
}

type RemoteConfiger struct {
	BaseConfiger
	RemoteProvider string
	RemoteHost     string
	ConfigKey      string
	ConfigType     string
}

func (rConfiger *RemoteConfiger) loadData() error {
	viper.AddRemoteProvider(rConfiger.RemoteProvider, rConfiger.RemoteHost, rConfiger.ConfigKey)
	viper.SetConfigType(rConfiger.ConfigType)
	return viper.ReadRemoteConfig()
}

func (rConfiger *RemoteConfiger) parse() error {
	var err = rConfiger.unmarshal()
	if err != nil {
		return err
	}

	if rConfiger.WatchChange {
		go func(rConfiger IConfiger) {
			for {
				time.Sleep(time.Second * 5) // delay after each request

				// currently, only tested with etcd support
				err := viper.WatchRemoteConfig()
				if err != nil {
					log.Errorf(context.Background(), "unable to read remote config: %v", err)
					continue
				}

				// unmarshal new config into our runtime config struct. you can also use channel
				// to implement a signal to notify the system of the changes
				rConfiger.unmarshal()
			}
		}(rConfiger)
	}

	return nil
}

// 1. For static config file, the configPath will be <path>/<file>. i.e. ./conf/app.yaml
// 2. For remote key/value store, the configPath will be <kvstoreType>#<kvstoreHost>#<kvstoreKey>#kvstoreType。
//    i.e. etcd|http://127.0.0.1:4001|/config/hugo.yaml|yaml. Currently, only etcd, consul, and Google filestore are supported
func ParseConfig(configPath string, rawVal interface{}, watchChange bool) error {
	var err error
	// First check whether it is remote config
	var sections []string = strings.Split(configPath, "|")
	if len(sections) > 0 {
		if len(sections) != 4 {
			return fmt.Errorf("it looks like a remote config, but the format %s is not correct", configPath)
		} else {
			var remoteProvider string = sections[0]
			if !stringInSlice(remoteProvider, viper.SupportedRemoteProviders) {
				return fmt.Errorf("remote provider %s is not supported yet", remoteProvider)
			}

			var configType string = sections[3]
			if !stringInSlice(configType, viper.SupportedExts) {
				return fmt.Errorf("config type %s is not supported yet", configType)
			}

			var remoteHost string = sections[1]
			var configKey string = sections[2]

			log.Infof(context.Background(), "TBD, RemoteProvider: %s, RemoteHost: %s, ConfigKey: %s, ConfigType: %s", remoteProvider, remoteHost, configKey, configType)
		}

	} else {
		var dir, file string = path.Split(configPath)
		if dir == "" {
			log.Infof(context.Background(), "default to current working folder .")
			dir = "."
		}

		if file == "" {
			log.Infof(context.Background(), "default to config file app.yaml")
			file = "app.yaml"
		}

		if _, err = os.Stat(path.Join(dir, file)); os.IsNotExist(err) {
			return fmt.Errorf("config file %s doesn't exist: %v", path.Join(dir, file), err)
		}

		var fileName, fileExt string
		var fileNameExt []string = strings.Split(file, ".")
		fileName = fileNameExt[0]
		fileExt = "yaml" // Default to "yaml"
		if len(fileNameExt) == 2 {
			fileExt = fileNameExt[1]
		}
		if !stringInSlice(fileExt, viper.SupportedExts) {
			return fmt.Errorf("config type %s is not supported yet", fileExt)
		}

		Configer = NewYamlConfiger(fileName, fileExt, dir, watchChange, rawVal)
		err = Configer.loadData()
		if err != nil {
			return fmt.Errorf("fail to load config file %s: %v", path.Join(dir, file), err)
		}

		err = Configer.parse()
		if err != nil {
			return fmt.Errorf("fail to unmarshal config file %s: %v", path.Join(dir, file), err)
		}
	}

	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
