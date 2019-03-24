package search

import (
	"github.com/joaosoft/dbr"
	"github.com/joaosoft/elastic"
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
	migration "github.com/joaosoft/migration/services"
)

type Search struct {
	maxSize       int
	config        *SearchConfig
	isLogExternal bool
	pm            *manager.Manager
	logger        logger.ILogger
}

type searchResult struct {
	Result     interface{} `json:"result"`
	Metadata   interface{} `json:"Metadata,omitempty"`
	Pagination *pagination `json:"pagination,omitempty"`
}

type pagination struct {
	First    *string `json:"first"`
	Previous *string `json:"previous"`
	Next     *string `json:"next"`
	Last     *string `json:"last"`
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

	// execute migrations
	if search.config.Migration != nil {
		migrationService, err := migration.NewCmdService(migration.WithCmdConfiguration(search.config.Migration))
		if err != nil {
			return nil, err
		}

		if _, err := migrationService.Execute(migration.OptionUp, 0, migration.ExecutorModeDatabase); err != nil {
			return nil, err
		}
	}

	return search, nil
}

func (search *Search) NewDatabaseSearch(stmt *dbr.StmtSelect) *searchHandler {
	return search.newSearchHandler(search.newDatabaseClient(stmt))
}

func (search *Search) NewElasticSearch(stmt *elastic.SearchService) *searchHandler {
	return search.newSearchHandler(search.newElasticClient(stmt))
}
