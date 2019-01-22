package search

import (
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type Search struct {
	config        *SearchConfig
	isLogExternal bool
	pm            *manager.Manager
	logger        logger.ILogger
}

// New ...
func New(options ...SearchOption) (*Search, error) {
	config, simpleConfig, err := NewConfig()

	search := &Search{
		pm:     manager.NewManager(manager.WithRunInBackground(true)),
		logger: logger.NewLogDefault("search", logger.WarnLevel),
		config: config.Search,
	}

	if search.isLogExternal {
		search.pm.Reconfigure(manager.WithLogger(search.logger))
	}

	if err != nil {
		search.logger.Error(err.Error())
	} else if config.Search != nil {
		search.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Search.Log.Level)
		search.logger.Debugf("setting log level to %s", level)
		search.logger.Reconfigure(logger.WithLevel(level))
	}

	search.Reconfigure(options...)

	return search, nil
}

func (search *Search) NewSearch() {
}
