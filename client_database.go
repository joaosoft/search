package search

import (
	"fmt"
	"reflect"

	"github.com/joaosoft/dbr"
)

type databaseClient struct {
	*dbr.StmtSelect
}

func (search *Search) newDatabaseClient(stmt *dbr.StmtSelect) *databaseClient {
	return &databaseClient{StmtSelect: stmt}
}

func (client *databaseClient) Exec(searchData *searchData) (int, error) {
	var err error

	// query
	for key, value := range searchData.query {
		client.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// search
	lenQ := len(searchData.searchFilters)
	if searchData.search != nil && lenQ > 0 {
		queryFilter := ""
		for i, filter := range searchData.searchFilters {
			queryFilter += fmt.Sprintf("%s ILIKE %s", filter, client.Db.Dialect.Encode("%"+*searchData.search+"%"))

			if i < lenQ-1 {
				queryFilter += " OR "
			}
		}

		client.Where(fmt.Sprintf("(%s)", queryFilter))
	}

	// pagination
	total := 0
	if searchData.hasPagination {
		_, err = client.Dbr.Select("count(1)").From(dbr.As(client.StmtSelect, "search")).Load(&total)

		if err != nil {
			return 0, err
		}

		if total == 0 {
			return 0, nil
		}
	}

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

	_, err = client.Load(searchData.object)
	if err != nil {
		return 0, err
	}

	// Metadata
	if searchData.hasMetadata {
		for _, item := range searchData.metadata {
			// function
			if item.function != nil {
				if err = item.function(reflect.ValueOf(searchData.object).Elem().Interface(), item.object, searchData.metadata); err != nil {
					return 0, err
				}
			}

			// statement
			if item.stmt != nil {
				if stmt, ok := item.stmt.(*dbr.StmtSelect); ok {
					stmt.Load(item.object)
				}
			}
		}
	}

	return total, err
}
