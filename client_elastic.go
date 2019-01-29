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
	// query
	terms := make([]elastic.Query, 0)
	for key, value := range searchData.query {
		terms = append(terms, elastic.NewTerm(key, value))
	}
	client.Query(elastic.NewBool().Must(terms...))

	// search
	lenQ := len(searchData.searchFilters)
	if searchData.search != nil && lenQ > 0 {
		queryString := elastic.NewQueryString(*searchData.search)
		for _, filter := range searchData.searchFilters {
			queryString.Fields(filter)
		}
		client.Query(queryString)
	}

	// pagination
	total := 0
	if searchData.hasPagination {
		response, err := client.Count()
		if err != nil {
			return 0, err
		}

		if response.OnError != nil || response.OnErrorDocumentNotFound != nil {
			return 0, nil
		}

		if response.Count == 0 {
			return 0, nil
		}

		total = int(response.Count)
	}

	if searchData.size > 0 {
		client.Size(searchData.size)
	}

	if searchData.page > 0 {
		client.From((searchData.page - 1) * searchData.size)
	}

	// order by
	sorts := make([]*elastic.SortField, 0)
	for _, order := range searchData.orders {
		switch order.direction {
		case orderAsc:
			sorts = append(sorts, elastic.NewSortField(order.column, elastic.OrderAsc))
		case orderDesc:
			sorts = append(sorts, elastic.NewSortField(order.column, elastic.OrderDesc))
		}
	}

	if _, err := client.Object(searchData.object).Query(elastic.NewSort(sorts...)).Search(); err != nil {
		return 0, err
	}

	// Metadata
	if searchData.hasMetadata {
		for _, item := range searchData.metadata {
			// function
			if item.function != nil {
				if err := item.function(reflect.ValueOf(searchData.object).Elem().Interface(), item.object, searchData.metadata); err != nil {
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
	return total, err
}
