package elastic

import (
	"encoding/json"

	"github.com/joaosoft/web"

	"fmt"

	"github.com/joaosoft/errors"
)

type IndexResponse struct {
	Acknowledged       bool   `json:"acknowledged"`
	ShardsAcknowledged bool   `json:"shards_acknowledged"`
	Index              string `json:"index"`
	*OnError
}

type IndexService struct {
	client *Elastic
	index  string
	typ    string
	body   []byte
}

func NewIndexService(e *Elastic) *IndexService {
	return &IndexService{
		client: e,
	}
}

func (e *IndexService) Index(index string) *IndexService {
	e.index = index
	return e
}

func (e *IndexService) Body(body interface{}) *IndexService {
	switch v := body.(type) {
	case []byte:
		e.body = v
	default:
		e.body, _ = json.Marshal(v)
	}
	return e
}

func (e *IndexService) Create() (*IndexResponse, error) {
	return e.execute(web.MethodPut)
}

func (e *IndexService) Update() (*IndexResponse, error) {
	return e.execute(web.MethodPut)
}

func (e *IndexService) Delete() (*IndexResponse, error) {
	return e.execute(web.MethodDelete)
}

func (e *IndexService) Exists() (bool, error) {
	_, err := e.execute(web.MethodHead)
	return err == nil, err
}

func (e *IndexService) execute(method web.Method) (*IndexResponse, error) {

	typ := ""
	if e.typ != "" {
		typ = fmt.Sprintf("/%s", e.typ)
	}

	request, err := e.client.NewRequest(method, fmt.Sprintf("%s/%s%s", e.client.config.Endpoint, e.index, typ), web.ContentTypeApplicationJSON, nil)
	if err != nil {
		return nil, errors.New(errors.LevelError, 0, err)
	}

	response, err := request.WithBody(e.body).Send()
	if err != nil {
		return nil, errors.New(errors.LevelError, 0, err)
	}

	// unmarshal data
	elasticResponse := IndexResponse{}

	if method != web.MethodHead {
		err = json.Unmarshal(response.Body, &elasticResponse)
		if err != nil {
			return nil, errors.New(errors.LevelError, 0, err)
		}
	}

	return &elasticResponse, nil
}
