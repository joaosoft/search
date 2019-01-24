package search

import (
	"fmt"
	"math"
	"strconv"
)

type searchHandler struct {
	client   searchClient
	path     string
	query    map[string]string
	search   *string
	filters  []string
	metadata map[string]*metadata
	orders   orders
	page     int
	size     int
	object   interface{}
}

type metadata struct {
	stmt   interface{}
	object interface{}
}

func newSearchHandler(client searchClient) *searchHandler {
	return &searchHandler{client: client, query: make(map[string]string), metadata: make(map[string]*metadata)}
}

func (searchHandler *searchHandler) Query(query map[string]string) *searchHandler {
	for key, value := range query {
		switch key {
		case constPage:
			searchHandler.page, _ = strconv.Atoi(value)
		case constSize:
			searchHandler.size, _ = strconv.Atoi(value)
		case constSearch:
			searchHandler.search = &value
		default:
			searchHandler.query[key] = value
		}
	}

	return searchHandler
}

func (searchHandler *searchHandler) Filters(fields ...string) *searchHandler {
	searchHandler.filters = fields
	return searchHandler
}

func (searchHandler *searchHandler) Metadata(name string, stmt interface{}, object interface{}) *searchHandler {
	searchHandler.metadata[name] = &metadata{stmt: stmt, object: object}
	return searchHandler
}

func (searchHandler *searchHandler) OrderBy(field string, direction direction) *searchHandler {
	searchHandler.orders = append(searchHandler.orders, &order{column: field, direction: direction})
	return searchHandler
}

func (searchHandler *searchHandler) Search(value string) *searchHandler {
	searchHandler.search = &value
	return searchHandler
}

func (searchHandler *searchHandler) Page(page int) *searchHandler {
	searchHandler.page = page
	return searchHandler
}

func (searchHandler *searchHandler) Path(path string) *searchHandler {
	searchHandler.path = path
	return searchHandler
}

func (searchHandler *searchHandler) Size(size int) *searchHandler {
	searchHandler.size = size
	return searchHandler
}

func (searchHandler *searchHandler) Bind(object interface{}) *searchHandler {
	searchHandler.object = object
	return searchHandler
}

func (searchHandler *searchHandler) Exec() (*searchResult, error) {
	searchData := &searchData{
		path:     searchHandler.path,
		query:    searchHandler.query,
		search:   searchHandler.search,
		filters:  searchHandler.filters,
		orders:   searchHandler.orders,
		page:     searchHandler.page,
		size:     searchHandler.size,
		object:   searchHandler.object,
		metadata: searchHandler.metadata,
	}
	total, err := searchHandler.client.Exec(searchData)

	// metadata
	metadata := make(map[string]interface{})
	for name, item := range searchData.metadata {
		metadata[name] = item.object
	}
	return &searchResult{
		Result:     searchHandler.object,
		Metadata:   metadata,
		Pagination: newPagination(searchData, total),
	}, err
}

func newPagination(searchData *searchData, total int) *pagination {
	pagination := pagination{}
	totalPages := int(math.Ceil(float64(total) / float64(searchData.size)))

	// if there are no results
	if total == 0 {
		return &pagination
	}

	// first page
	if totalPages > 1 && searchData.page > 1 {
		first := fmt.Sprintf("%s?page=%d&size=%d", searchData.path, 1, searchData.size)
		pagination.First = &first

		// previous page
		previous := fmt.Sprintf("%s?page=%d&size=%d", searchData.path, searchData.page-1, searchData.size)
		pagination.Previous = &previous
	}

	// next page
	if totalPages > searchData.page {
		next := fmt.Sprintf("%s?page=%d&size=%d", searchData.path, searchData.page+1, searchData.size)
		pagination.Next = &next

		// last page
		size := searchData.size

		// calculate the remainder items on the last page
		//remainder := total % searchData.size
		//if remainder > 0 {
		//	size = remainder
		//}

		last := fmt.Sprintf("%s?page=%d&size=%d", searchData.path, totalPages, size)
		pagination.Last = &last
	}

	return &pagination
}
