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

func (client *DatabaseClient) Exec(searchData *searchData) (int, error) {
	var err error

	// query
	for key, value := range searchData.query {
		client.Where(key, value)
	}

	// filters and search
	for _, filter := range searchData.filters {
		if searchData.search != nil {
			client.Where(fmt.Sprintf("%s ILIKE %s", filter, *searchData.search))
		}
	}

	total := 0
	_, err = client.Dbr.Select("count(1)").From(dbr.As(client.StmtSelect, "search")).Load(&total)

	if err != nil {
		return 0, err
	}

	if total > 0 {
		// pagination
		if searchData.size > 0 {
			client.Limit(searchData.size)
		}

		if searchData.page > 0 {
			client.Offset((searchData.page - 1) * searchData.size)
		}

		// order by
		for _, order := range searchData.orders {
			switch order.direction {
			case orderAsc:
				client.OrderAsc(order.column)
			case orderDesc:
				client.OrderDesc(order.column)
			}
		}

		_, err := client.Load(searchData.object)

		// load metadata
		for _, item := range searchData.metadata {
			if stmt, ok := item.stmt.(*dbr.StmtSelect); ok {
				stmt.Load(item.object)
			}
		}
		return total, err
	}

	return 0, err
}
