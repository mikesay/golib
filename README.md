# golib
A common go library

## Config module
Wrap Viper(https://github.com/spf13/viper) to provide simple interface for Go application.

## Test yaml
```yaml
---
YamlConfig:
  Log:
    # Levels(1:fatal 2:error,3:warn,4:info,5:debug,6:trace)
    Level: 5
    # Level format(support types：text/json)
    Format: "json"
    # Output(support types: stdout/stderr)
    Output: "stdout"
  System:
    scrapeInterval: 600
  AliSLS:
    - host: cn-shanghai.log.aliyuncs.com
      ackId: xxxx
      ackSec: xxxx
      logProjects:
        - projectName: mikesay-test-project
          logStores:
            - logStoreName: istiolog
              collectors:
                - istio
                - istio1
            - logStoreName: istiolog1
              collectors:
                - istio
```

```go
// SystemSetting
type SystemSetting struct {
	ScrapeInterval int64 `mapstructure:"scrapeInterval"`
}

// Log settings
type LogSetting struct {
	Level  int    `mapstructure:"Level"`
	Format string `mapstructure:"Format"`
	Output string `mapstructure:"Output"`
}

// AliCloud SLS settings
type AliSLSLogStore struct {
	LogstoreName string   `mapstructure:"logStoreName"`
	Collectors   []string `mapstructure:"collectors"`
}

type AliSLSLogProject struct {
	ProjectName string           `mapstructure:"projectName"`
	LogStores   []AliSLSLogStore `mapstructure:"logStores"`
}

type AliSLSSetting struct {
	Host        string             `mapstructure:"host"`
	AckID       string             `mapstructure:"ackId"`
	AckSec      string             `mapstructure:"ackSec"`
	LogProjects []AliSLSLogProject `mapstructure:"logProjects"`
}

type AliSLSSettings []AliSLSSetting

// Splunk settings
type SplunkSourceLogStores struct {
	LogstoreName string   `mapstructure:"logStoreName"`
	Collectors   []string `mapstructure:"collectors"`
}

type SplunkIndex struct {
	IndexName       string                  `mapstructure:"indexName"`
	SourceLogStores []SplunkSourceLogStores `mapstructure:"sourceLogStores"`
}

type SplunkSetting struct {
	Host     string        `mapstructure:"host"`
	Username string        `mapstructure:"username"`
	Password string        `mapstructure:"password"`
	Indexes  []SplunkIndex `mapstructure:"indexes"`
}

type SplunkSettings []SplunkSetting

type YamlConfigSetting struct {
	SysSetting    SystemSetting  `mapstructure:"System"`
	LogSetting    LogSetting     `mapstructure:"Log"`
	AliSLSSetting AliSLSSettings `mapstructure:"AliSLS"`
}

type ConfigS struct {
	Ycs YamlConfigSetting `mapstructure:"YamlConfig"`
}

func main() {
	appcfg := AppConfig{}
	err := config.ParseConfig("./conf/app.ini", &appcfg, true)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshal config: %w", err))
	}

	go func(cfg *AppConfig) {
		for {
			fmt.Printf("Httpport is %s\n", cfg.Default.Appname)
			time.Sleep(time.Second * 2)
		}
	}(&appcfg)
}
```

## Test ini
```ini
appname = beeapi
httpport = 8080
runmode = dev
autorender = false
copyrequestbody = true
EnableDocs = true
sqlconn = 

[dev]
httpport = 8081
```

```go
type DevSection struct {
	Httpport string `mapstructure:"httpport"`
}

type DefaultSection struct {
	Appname         string `mapstructure:"appname"`
	Httpport        string `mapstructure:"httpport"`
	Runmode         string `mapstructure:"runmode"`
	Autorender      bool   `mapstructure:"autorender"`
	Copyrequestbody bool   `mapstructure:"copyrequestbody"`
	EnableDocs      bool   `mapstructure:"EnableDocs"`
	Sqlconn         string `mapstructure:"sqlconn"`
}

type AppConfig struct {
	Default DefaultSection `mapstructure:"default"`
	Dev     DevSection     `mapstructure:"dev"`
}

func main() {
	appcfg := AppConfig{}
	err := config.ParseConfig("./conf/app.ini", &appcfg, true)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshal config: %w", err))
	}

	go func(cfg *AppConfig) {
		for {
			fmt.Printf("Httpport is %s\n", cfg.Default.Appname)
			time.Sleep(time.Second * 2)
		}
	}(&appcfg)
}
```
