package search

import (
	logger "github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// SearchOption ...
type SearchOption func(search *Search)

// Reconfigure ...
func (search *Search) Reconfigure(options ...SearchOption) {
	for _, option := range options {
		option(search)
	}
}

// WithConfiguration ...
func WithConfiguration(config *SearchConfig) SearchOption {
	return func(search *Search) {
		search.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) SearchOption {
	return func(search *Search) {
		search.logger = logger
		search.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) SearchOption {
	return func(search *Search) {
		search.logger.SetLevel(level)
	}
}

// WithManager ...
func WithManager(mgr *manager.Manager) SearchOption {
	return func(search *Search) {
		search.pm = mgr
	}
}

// WithMaxSize ...
func WithMaxSize(maxSize int) SearchOption {
	return func(search *Search) {
		search.maxSize = maxSize
	}
}
