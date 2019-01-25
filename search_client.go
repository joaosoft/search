package search

type searchClient interface {
	Exec(searchData *searchData) (int, error)
}

type searchData struct {
	hasPagination bool
	hasMetadata   bool
	path          string
	query         map[string]string
	search        *string
	filters       map[string]string
	searchFilters []string
	orders        orders
	page          int
	size          int
	object        interface{}
	metadata      map[string]*metadata
}
