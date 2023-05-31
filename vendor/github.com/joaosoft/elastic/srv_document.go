package elastic

import (
	"encoding/json"
	"github.com/joaosoft/web"

	"fmt"
	"github.com/joaosoft/errors"
)

type DocumentResponse struct {
	Index   string `json:"_index"`
	Type    string `json:"_type"`
	ID      string `json:"_id"`
	Version int64  `json:"_version"`
	Result  string `json:"result"`
	Shards  struct {
		Total      int64 `json:"total"`
		Successful int64 `json:"successful"`
		Failed     int64 `json:"failed"`
	} `json:"_shards"`
	SeqNo       int64 `json:"_seq_no"`
	PrimaryTerm int64 `json:"_primary_term"`
	OnError
}

type DocumentService struct {
	client     *Elastic
	index      string
	typ        string
	id         string
	body       []byte
	parameters map[string]interface{}
}

type RefreshType string

const (
	Refresh   RefreshType = "true"
	WaitFor   RefreshType = "wait_for"
	NoRefresh RefreshType = "false"
)

func NewDocumentService(e *Elastic) *DocumentService {
	return &DocumentService{
		client:     e,
		parameters: make(map[string]interface{}),
	}
}

func (e *DocumentService) Index(index string) *DocumentService {
	e.index = index
	return e
}

func (e *DocumentService) Type(typ string) *DocumentService {
	e.typ = typ
	return e
}

func (e *DocumentService) Id(id string) *DocumentService {
	e.id = id
	return e
}

func (e *DocumentService) Refresh(refreshType RefreshType) *DocumentService {
	e.parameters["refresh"] = refreshType
	return e
}

func (e *DocumentService) Body(body interface{}) *DocumentService {
	switch v := body.(type) {
	case []byte:
		e.body = v
	default:
		e.body, _ = json.Marshal(v)
	}
	return e
}

func (e *DocumentService) Create() (*DocumentResponse, error) {
	return e.execute(web.MethodPost)
}

func (e *DocumentService) Update() (*DocumentResponse, error) {
	return e.execute(web.MethodPut)
}

func (e *DocumentService) Delete() (*DocumentResponse, error) {
	return e.execute(web.MethodDelete)
}

func (e *DocumentService) execute(method web.Method) (*DocumentResponse, error) {

	var query string

	if e.id != "" {
		query += fmt.Sprintf("/%s", e.id)
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

	request, err := e.client.Client.NewRequest(method, fmt.Sprintf("%s/%s/%s%s", e.client.config.Endpoint, e.index, e.typ, query), web.ContentTypeApplicationJSON, nil)
	if err != nil {
		return nil, errors.New(errors.LevelError, 0, err)
	}

	response, err := request.WithBody(e.body).Send()
	if err != nil {
		return nil, errors.New(errors.LevelError, 0, err)
	}

	elasticResponse := DocumentResponse{}
	if err = json.Unmarshal(response.Body, &elasticResponse); err != nil {
		return nil, errors.New(errors.LevelError, 0, err)
	}

	return &elasticResponse, nil
}
