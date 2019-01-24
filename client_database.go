package search

import (
	"fmt"

	"github.com/joaosoft/dbr"
)

type DatabaseClient struct {
	*dbr.StmtSelect
}

func newDatabaseClient(stmt *dbr.StmtSelect) *DatabaseClient {
	return &DatabaseClient{StmtSelect: stmt}
}

func (client *DatabaseClient) Exec(query map[string]string, search *string, filters []string, orders orders, page int, size int, object interface{}) (int, error) {

	// query
	for key, value := range query {
		client.Where(key, value)
	}

	// filters and search
	for _, filter := range filters {
		if search != nil {
			client.Where(fmt.Sprintf("%s ILIKE %s", filter, fmt.Sprintf("%%s%", search)))
		}
	}

	stmtSelect, err := client.Build()
	if err != nil {
		return 0, err
	}

	client.Dbr.Select("count(1)").From(dbr.Field(stmtSelect).As("search"))

	// pagination
	if size > 0 {
		client.Limit(size)
	}

	if page > 0 {
		client.Offset(page * size)
	}

	// order by
	for _, order := range orders {
		switch order.direction {
		case orderAsc:
			client.OrderAsc(order.column)
		case orderDesc:
			client.OrderDesc(order.column)
		}
	}

	return client.Load(object)
}
