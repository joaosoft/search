package elastic

import (
	"sync"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
	"github.com/joaosoft/web"
)

type Elastic struct {
	*web.Client
	config        *ElasticConfig
	isLogExternal bool
	logger        logger.ILogger
	pm            *manager.Manager
	mux           sync.Mutex
}

// NewElastic ...
func NewElastic(options ...ElasticOption) (*Elastic, error) {
	config, simpleConfig, err := NewConfig()
	webClient, err := web.NewClient()
	if err != nil {
		return nil, err
	}

	service := &Elastic{
		Client: webClient,
		pm:     manager.NewManager(manager.WithRunInBackground(false)),
		config: config.Elastic,
		logger: logger.NewLogDefault("elastic", logger.WarnLevel),
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else if config.Elastic != nil {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Elastic.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(service.logger))
	}

	service.Reconfigure(options...)

	return service, nil
}

func (e *Elastic) Document() *DocumentService {
	return NewDocumentService(e)
}

func (e *Elastic) Search() *SearchService {
	return NewSearchService(e)
}

func (e *Elastic) Index() *IndexService {
	return NewIndexService(e)
}

func (e *Elastic) Bulk() *BulkService {
	return NewBulkService(e)
}
