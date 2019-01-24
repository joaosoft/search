package search

import "github.com/joaosoft/elastic"

type ElasticClient struct {
	elastic.Elastic
}

func (client *ElasticClient) Exec(query map[string]string, search *string, filters []string, orders orders, page int, size int, object interface{}) (int, error) {

	return 0, client.Search().Query("").Object(object).Execute()
}
