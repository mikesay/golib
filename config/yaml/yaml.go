package yaml

import (
	"github.com/mikesay/golib/config"
	"github.com/spf13/viper"
)

type YamlConfiger struct{}

func (yamlConfiger *YamlConfiger) Unmarshal(rawVal interface{}) error {
	return viper.Unmarshal(rawVal)
}

func init() {
	config.Register("yaml", new(YamlConfiger))
}
