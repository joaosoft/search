package search

import (
	"strconv"

	"github.com/joaosoft/dbr"
)

type searchHandler struct {
	client     searchClient
	selectStmt *dbr.StmtSelect
	query      map[string]string
	search     *string
	filters    []string
	orders     orders
	page       int
	size       int
	object     interface{}
}

func newSearchHandler(client searchClient) *searchHandler {
	return &searchHandler{client: client, query: make(map[string]string)}
}

func (searchHandler *searchHandler) Builder(selectStmt *dbr.StmtSelect) *searchHandler {
	searchHandler.selectStmt = selectStmt
	return searchHandler
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

func (searchHandler *searchHandler) Size(size int) *searchHandler {
	searchHandler.size = size
	return searchHandler
}

func (searchHandler *searchHandler) Bind(object interface{}) *searchHandler {
	searchHandler.object = object
	return searchHandler
}

func (searchHandler *searchHandler) Exec() (*searchResult, error) {
	err := searchHandler.client.
		Exec(searchHandler.selectStmt,
			searchHandler.query,
			searchHandler.search,
			searchHandler.filters,
			searchHandler.orders,
			searchHandler.page,
			searchHandler.size,
			searchHandler.object)

	return &searchResult{
		data:       searchHandler.object,
		pagination: &pagination{},
	}, err
}
