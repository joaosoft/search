package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/joaosoft/errors"
	"github.com/joaosoft/web"
)

type operation string

const (
	constScroll = "scroll"
	constFrom   = "from"
	constSize   = "size"

	constOperationSearch = "_search"
	constOperationCount  = "_count"
)

type Query interface {
	Data() interface{}
}

type OnCount struct {
	Count int64 `json:"count"`
}

type SearchResponse struct {
	Took     int64 `json:"took"`
	TimedOut bool  `json:"timed_out"`
	Shards   struct {
		Total      int64 `json:"total"`
		Successful int64 `json:"successful"`
		Skipped    int64 `json:"skipped"`
		Failed     int64 `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total    int64   `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string          `json:"_index"`
			Type   string          `json:"_type"`
			ID     string          `json:"_id"`
			Score  float64         `json:"_score"`
			Source json.RawMessage `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
	*OnCount
	*OnError
	*OnErrorDocumentNotFound
}

type OnErrorDocumentNotFound struct {
	Index string `json:"_index"`
	Type  string `json:"_type"`
	ID    string `json:"_id"`
	Found bool   `json:"found"`
}

type SearchService struct {
	client     *Elastic
	index      []string
	typ        []string
	id         string
	queries    map[string]interface{}
	template   []byte
	query      map[string]interface{}
	body       []byte
	object     interface{}
	parameters map[string]interface{}
	method     web.Method
	operation  operation
}

func NewSearchService(client *Elastic) *SearchService {
	return &SearchService{
		client:     client,
		method:     web.MethodGet,
		parameters: make(map[string]interface{}),
		queries:    make(map[string]interface{}),
		query:      make(map[string]interface{}),
		operation:  "_search",
	}
}

func (e *SearchService) Index(index ...string) *SearchService {
	e.index = index
	return e
}

func (e *SearchService) Type(typ ...string) *SearchService {
	e.typ = typ
	return e
}

func (e *SearchService) Id(id string) *SearchService {
	e.id = id
	return e
}

func (e *SearchService) Body(body []byte) *SearchService {
	e.body = body
	return e
}

func (e *SearchService) Data() *SearchService {
	if len(e.index) > 1 || len(e.typ) > 1 {
		e.queries["index"] = strings.Join(e.index, ",")
		e.queries["type"] = strings.Join(e.typ, ",")
	}
	return e
}
func (e *SearchService) Query(queries ...Query) *SearchService {

	for _, query := range queries {
		for key, value := range query.Data().(map[string]interface{}) {
			e.queries[key] = value
		}
	}

	return e
}

func (e *SearchService) Object(object interface{}) *SearchService {
	e.object = object
	return e
}

func (e *SearchService) From(from int) *SearchService {
	e.parameters[constFrom] = from
	return e
}

func (e *SearchService) Size(size int) *SearchService {
	e.parameters[constSize] = size
	return e
}

func (e *SearchService) Scroll(scrollTime string) *SearchService {
	e.parameters[constScroll] = scrollTime
	return e
}

type SearchTemplate struct {
	Data interface{} `json:"data,omitempty"`
	From int         `json:"from,omitempty"`
	Size int         `json:"size,omitempty"`
}

type CountTemplate struct {
	Data interface{} `json:"data,omitempty"`
}

func (e *SearchService) Template(path, name string, data interface{}, reload bool) *SearchService {
	key := fmt.Sprintf("%s/%s", path, name)

	var result bytes.Buffer
	var err error

	if _, found := templates[key]; !found {
		e.client.mux.Lock()
		defer e.client.mux.Unlock()
		templates[key], err = ReadFile(key, nil)
		if err != nil {
			e.client.logger.Error(err)
			return e
		}
	}

	t := template.New(name)
	t, err = t.Parse(string(templates[key]))
	if err == nil {
		if err := t.ExecuteTemplate(&result, name, data); err != nil {
			e.client.logger.Error(err)
			return e
		}

		e.template = result.Bytes()
	} else {
		e.client.logger.Error(err)
		return e
	}

	return e
}

func (e *SearchService) Count() (*SearchResponse, error) {
	e.operation = constOperationCount
	return e.execute()
}

func (e *SearchService) Search() (*SearchResponse, error) {
	e.operation = constOperationSearch
	return e.execute()
}

func (e *SearchService) execute() (*SearchResponse, error) {

	if len(e.queries) > 0 || len(e.index) > 1 || len(e.typ) > 1 {
		e.query["query"] = e.Data()
		e.body, _ = json.Marshal(e.query)
	} else if e.template != nil {
		e.body = e.template
	}

	if e.body != nil {
		e.method = web.MethodPost
	}

	var query string
	if e.id != "" {
		query += fmt.Sprintf("/%s", e.id)
	} else {
		query += fmt.Sprintf("/%s", e.operation)
	}

	lenQ := len(e.parameters)
	if lenQ > 0 {
		query += "?"
	}

	addSeparator := false
	for name, value := range e.parameters {
		if addSeparator {
			query += "&"
		}

		query += fmt.Sprintf("%s=%+v", name, value)
		addSeparator = true
	}

	request, err := e.client.Client.NewRequest(e.method, fmt.Sprintf("%s/%s%s", e.client.config.Endpoint, e.index[0], query), web.ContentTypeApplicationJSON, nil)
	if err != nil {
		return nil, errors.New(errors.LevelError, 0, err)
	}

	response, err := request.WithBody(e.body).Send()
	if err != nil {
		return nil, errors.New(errors.LevelError, 0, err)
	}

	elasticResponse := SearchResponse{}
	if err := json.Unmarshal(response.Body, &elasticResponse); err != nil {
		e.client.logger.Error(err)
		return nil, errors.New(errors.LevelError, 0, err)
	}

	if elasticResponse.OnError != nil {
		return &elasticResponse, nil
	}

	if e.operation == constOperationSearch {
		rawHits := make([]json.RawMessage, len(elasticResponse.Hits.Hits))
		for i, rawHit := range elasticResponse.Hits.Hits {
			rawHits[i] = rawHit.Source
		}

		arrayHits, err := json.Marshal(rawHits)
		if err != nil {
			return nil, errors.New(errors.LevelError, 0, err)
		}

		if err := json.Unmarshal(arrayHits, e.object); err != nil {
			return nil, errors.New(errors.LevelError, 0, err)
		}
	}

	return &elasticResponse, nil
}
