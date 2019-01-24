package search

type searchClient interface {
	Exec(query map[string]string, search *string, filters []string, orders orders, page int, size int, object interface{}) (int, error)
}
