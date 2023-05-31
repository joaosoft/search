package elastic

import (
	"fmt"

	"github.com/joaosoft/logger"

	"github.com/joaosoft/manager"
)

// AppConfig ...
type AppConfig struct {
	Elastic *ElasticConfig `json:"elastic"`
}

// ElasticConfig ...
type ElasticConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
	Endpoint string `json:"endpoint"`
}

// NewConfig ...
func NewConfig(config ...interface{}) (*AppConfig, manager.IConfig, error) {
	appConfig := &AppConfig{}
	simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	if len(config) > 0 {

		if appConfig.Elastic == nil {
			appConfig.Elastic = &ElasticConfig{}
			appConfig.Elastic.Log.Level = logger.ErrorLevel.String()
		}

		appConfig.Elastic.Endpoint = config[0].(string)
	}

	return appConfig, simpleConfig, err
}
