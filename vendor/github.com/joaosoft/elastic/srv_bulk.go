package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/joaosoft/errors"
	"github.com/joaosoft/web"
)

const (
	BulkIndexHeader  = `{ "index" : { "_index" : "%s", "_type" : "%s", "_id" : "%s" } }`
	BulkCreateHeader = `{ "create" : { "_index" : "%s", "_type" : "%s", "_id" : "%s" } }`
	BulkUpdateHeader = `{ "update" : { "_index" : "%s", "_type" : "%s", "_id" : "%s" } }`
	BulkDeleteHeader = `{ "delete" : { "_index" : "%s", "_type" : "%s", "_id" : "%s" } }`
)

type BulkResponse struct {
	Took   int64 `json:"took"`
	Errors bool  `json:"errors"`
	Items  []struct {
		Index struct {
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
			Status      int64 `json:"status"`
			SeqNo       int64 `json:"_seq_no"`
			PrimaryTerm int64 `json:"_primary_term"`
			*OnErrorBulkOperation
		} `json:"index,omitempty"`
		Delete struct {
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
			Status      int64 `json:"status"`
			SeqNo       int64 `json:"_seq_no"`
			PrimaryTerm int64 `json:"_primary_term"`
			*OnErrorBulkOperation
		} `json:"delete,omitempty"`
		Create struct {
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
			Status      int64 `json:"status"`
			SeqNo       int64 `json:"_seq_no"`
			PrimaryTerm int64 `json:"_primary_term"`
			*OnErrorBulkOperation
		} `json:"create,omitempty"`
		Update struct {
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
			Status      int64 `json:"status"`
			SeqNo       int64 `json:"_seq_no"`
			PrimaryTerm int64 `json:"_primary_term"`
			*OnErrorBulkOperation
		} `json:"update,omitempty"`
	} `json:"items"`
}
type BulkService struct {
	client *Elastic
	index  string
	typ    string
	id     string
	body   []byte
	method web.Method
	buffer *bytes.Buffer
}

func NewBulkService(e *Elastic) *BulkService {
	return &BulkService{
		buffer: bytes.NewBufferString(""),
		client: e,
		method: web.MethodPost,
	}
}

func (e *BulkService) Index(index string) *BulkService {
	e.index = index
	return e
}

func (e *BulkService) Type(typ string) *BulkService {
	e.typ = typ
	return e
}

func (e *BulkService) Id(id string) *BulkService {
	e.id = id
	return e
}

func (e *BulkService) Body(body interface{}) *BulkService {
	switch v := body.(type) {
	case []byte:
		e.body = v
	default:
		e.body, _ = json.Marshal(v)
	}
	return e
}

func (e *BulkService) DoIndex() error {
	e.buffer.WriteString(fmt.Sprintf(BulkIndexHeader, e.index, e.typ, e.id))
	e.buffer.WriteString("\n")
	e.buffer.Write(e.body)
	e.buffer.WriteString("\n")

	return nil
}

func (e *BulkService) DoCreate() error {
	e.buffer.WriteString(fmt.Sprintf(BulkCreateHeader, e.index, e.typ, e.id))
	e.buffer.WriteString("\n")
	e.buffer.Write(e.body)
	e.buffer.WriteString("\n")

	return nil
}

func (e *BulkService) DoUpdate() error {
	e.buffer.WriteString(fmt.Sprintf(BulkUpdateHeader, e.index, e.typ, e.id))
	e.buffer.WriteString("\n")
	e.buffer.Write(e.body)
	e.buffer.WriteString("\n")

	return nil
}

func (e *BulkService) DoDelete() error {
	e.buffer.WriteString(fmt.Sprintf(BulkDeleteHeader, e.index, e.typ, e.id))
	e.buffer.WriteString("\n")

	return nil
}

func (e *BulkService) Execute() (*BulkResponse, error) {

	if e.buffer.Len() == 0 {
		return nil, nil
	}

	e.buffer.WriteString("\n") // needs a blank line at the end

	request, err := e.client.Client.NewRequest(e.method, fmt.Sprintf("%s/_bulk", e.client.config.Endpoint), web.ContentTypeApplicationJSON, nil)
	if err != nil {
		return nil, errors.New(errors.LevelError, 0, err)
	}

	response, err := request.WithBody(e.buffer.Bytes()).Send()
	if err != nil {
		return nil, errors.New(errors.LevelError, 0, err)
	}

	elasticResponse := BulkResponse{}
	if err = json.Unmarshal(response.Body, &elasticResponse); err != nil {
		return nil, errors.New(errors.LevelError, 0, err)
	}

	if elasticResponse.Errors {
		return &elasticResponse, errors.New(errors.LevelError, 0, "error executing the request")
	}

	return &elasticResponse, nil
}
