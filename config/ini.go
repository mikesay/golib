package config

type IniConfiger struct {
	FileConfiger
}

func NewIniConfiger(configFile string,
	configName string,
	configType string,
	configPath string,
	watchChange bool,
	rawVal interface{}) IConfiger {

	var iniConfiger *IniConfiger = new(IniConfiger)
	iniConfiger.WatchChange = watchChange
	iniConfiger.ConfigName = configName
	iniConfiger.ConfigType = configType
	iniConfiger.ConfigPath = configPath
	iniConfiger.ConfigFile = configFile
	iniConfiger.RawVal = rawVal

	return iniConfiger
}

func init() {
	FileConfigerFactory["ini"] = NewIniConfiger
}
