package search

import "github.com/joaosoft/elastic"

type ElasticClient struct {
	elastic.Elastic
}

func (client *ElasticClient) Exec(builder builder, search *string, query map[string]string, page int, size int, object interface{}) error {
	str, err := builder.Build()
	if err != nil {
		return err
	}

	return client.Search().Query(str).Object(object).Execute()
}
