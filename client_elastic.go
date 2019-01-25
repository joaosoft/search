package search

import "github.com/joaosoft/elastic"

type elasticClient struct {
	*elastic.SearchService
}

func (search *Search) newElasticClient(stmt *elastic.SearchService) *elasticClient {
	return &elasticClient{SearchService: stmt}
}

func (client *elasticClient) Exec(searchData *searchData) (int, error) {
	return 0, client.Query("").Object(searchData.object).Execute()
}
