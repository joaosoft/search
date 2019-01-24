package search

import (
	"fmt"

	"github.com/joaosoft/dbr"
)

type DatabaseClient struct{}

func (client *DatabaseClient) Exec(selectStmt *dbr.StmtSelect, query map[string]string, search *string, filters []string, orders orders, page int, size int, object interface{}) (int, error) {

	// query
	for key, value := range query {
		selectStmt.Where(key, value)
	}

	// filters and search
	for _, filter := range filters {
		if search != nil {
			selectStmt.Where(fmt.Sprintf("%s ILIKE %s", filter, fmt.Sprintf("%%s%", search)))
		}
	}

	// order by
	for _, order := range orders {
		switch order.direction {
		case orderAsc:
			selectStmt.OrderAsc(order.column)
		case orderDesc:
			selectStmt.OrderDesc(order.column)
		}
	}

	return selectStmt.Load(object)
}
