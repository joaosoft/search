package search

import (
	"fmt"
	"html"
	"math"
	"strconv"
)

type searchHandler struct {
	client        searchClient
	hasPagination bool
	hasMetadata   bool
	path          string
	query         map[string]string
	search        *string
	filters       map[string]string
	searchFilters []string
	metadata      map[string]*metadata
	orders        orders
	page          int
	size          int
	maxSize       int
	object        interface{}
}

type metadataFunction func(result interface{}, object interface{}) error

type metadata struct {
	stmt     interface{}
	function metadataFunction
	object   interface{}
}

func (search *Search) newSearchHandler(client searchClient) *searchHandler {
	return &searchHandler{
		client:        client,
		query:         make(map[string]string),
		filters:       make(map[string]string),
		searchFilters: make([]string, 0),
		metadata:      make(map[string]*metadata),
		hasPagination: true,
		hasMetadata:   true,
	}
}

func (searchHandler *searchHandler) Query(query map[string]string) *searchHandler {
	for key, value := range query {
		value = html.UnescapeString(value)

		switch key {
		case constPage:
			searchHandler.page, _ = strconv.Atoi(value)
		case constSize:
			searchHandler.size, _ = strconv.Atoi(value)
		case constSearch:
			searchHandler.search = &value
		default:
			if filter, ok := searchHandler.filters[key]; ok {
				searchHandler.query[key] = filter
			}
		}
	}

	return searchHandler
}

func (searchHandler *searchHandler) Filters(fields ...string) *searchHandler {
	for _, field := range fields {
		searchHandler.filters[field] = field
	}
	return searchHandler
}

func (searchHandler *searchHandler) Filter(searchName string, internalName string) *searchHandler {
	searchHandler.filters[searchName] = internalName
	return searchHandler
}

func (searchHandler *searchHandler) SearchFilters(fields ...string) *searchHandler {
	searchHandler.searchFilters = append(searchHandler.searchFilters, fields...)
	return searchHandler
}

func (searchHandler *searchHandler) WithoutPagination() *searchHandler {
	searchHandler.hasPagination = false
	return searchHandler
}

func (searchHandler *searchHandler) WithoutMetadata() *searchHandler {
	searchHandler.hasMetadata = false
	return searchHandler
}

func (searchHandler *searchHandler) Metadata(name string, stmt interface{}, object interface{}) *searchHandler {
	searchHandler.metadata[name] = &metadata{stmt: stmt, object: object}
	return searchHandler
}

func (searchHandler *searchHandler) MetadataFunction(name string, function metadataFunction, object interface{}) *searchHandler {
	searchHandler.metadata[name] = &metadata{function: function, object: object}
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

func (searchHandler *searchHandler) MaxSize(maxSize int) *searchHandler {
	searchHandler.maxSize = maxSize
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

	if searchHandler.maxSize > 0 && searchHandler.size > searchHandler.maxSize {
		searchHandler.size = searchHandler.maxSize
	}

	searchData := &searchData{
		hasPagination: searchHandler.hasPagination,
		hasMetadata:   searchHandler.hasMetadata,
		path:          searchHandler.path,
		query:         searchHandler.query,
		search:        searchHandler.search,
		filters:       searchHandler.filters,
		searchFilters: searchHandler.searchFilters,
		orders:        searchHandler.orders,
		page:          searchHandler.page,
		size:          searchHandler.size,
		object:        searchHandler.object,
		metadata:      searchHandler.metadata,
	}
	total, err := searchHandler.client.Exec(searchData)

	// metadata
	var metadata map[string]interface{}
	if searchHandler.hasMetadata {
		metadata = make(map[string]interface{})
		for name, item := range searchData.metadata {
			metadata[name] = item.object
		}
	}

	// pagination
	var pagination *pagination
	if searchHandler.hasPagination {
		pagination = newPagination(searchData, total)
	}

	// result
	return &searchResult{
		Result:     searchHandler.object,
		Metadata:   metadata,
		Pagination: pagination,
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
