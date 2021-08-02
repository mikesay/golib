package config

type YamlConfiger struct {
	FileConfiger
}

func NewYamlConfiger(configName string,
	configType string,
	configPath string,
	watchChange bool,
	rawVal interface{}) IConfiger {

	var yamlConfiger *YamlConfiger = new(YamlConfiger)
	yamlConfiger.WatchChange = watchChange
	yamlConfiger.ConfigName = configName
	yamlConfiger.ConfigType = configType
	yamlConfiger.ConfigPath = configPath
	yamlConfiger.RawVal = rawVal

	return yamlConfiger
}

func init() {
	FileConfigerFactory["yaml"] = NewYamlConfiger
}
