package search

import "github.com/joaosoft/elastic"

type ElasticClient struct {
	elastic.Elastic
}

func (client *ElasticClient) Exec(searchData *searchData) (int, error) {

	return 0, client.Search().Query("").Object(searchData.object).Execute()
}
