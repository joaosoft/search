package search

import "github.com/joaosoft/dbr"

type searchClient interface {
	Exec(selectStmt *dbr.StmtSelect, query map[string]string, search *string, filters []string, orders orders, page int, size int, object interface{}) error
}
