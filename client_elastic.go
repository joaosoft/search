package search

import (
	"reflect"

	"github.com/joaosoft/elastic"
)

type elasticClient struct {
	*elastic.SearchService
}

func (search *Search) newElasticClient(stmt *elastic.SearchService) *elasticClient {
	return &elasticClient{SearchService: stmt}
}

func (client *elasticClient) Exec(searchData *searchData) (int, error) {
	// pagination
	if searchData.hasPagination {
		response, err := client.Count()
		if err != nil {
			return 0, err
		}

		if response.Count == 0 {
			return 0, nil
		}
	}

	if searchData.size > 0 {
		client.Size(searchData.size)
	}

	if searchData.page > 0 {
		client.From((searchData.page - 1) * searchData.size)
	}

	// order by
	for _, order := range searchData.orders {
		switch order.direction {
		case orderAsc:
			//client.OrderAsc(order.column)
		case orderDesc:
			//client.OrderDesc(order.column)
		}
	}

	if _, err := client.Object(searchData.object).Search(); err != nil {
		return 0, err
	}

	// metadata
	if searchData.hasMetadata {
		for _, item := range searchData.metadata {
			// function
			if item.function != nil {
				if err := item.function(reflect.ValueOf(searchData.object).Elem().Interface(), item.object); err != nil {
					return 0, err
				}
			}

			// statement
			if item.stmt != nil {
				if stmt, ok := item.stmt.(*elastic.SearchService); ok {
					if _, err := stmt.Object(item.object).Search(); err != nil {
						return 0, err
					}
				}
			}
		}
	}

	_, err := client.Object(searchData.object).Search()
	return 0, err
}
