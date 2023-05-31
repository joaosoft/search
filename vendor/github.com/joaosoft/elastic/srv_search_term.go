package elastic

import "encoding/json"

type Term struct {
	mappings map[string]interface{}
	field    string
}

func NewTerm(field string, value interface{}) *Term {
	new := &Term{
		mappings: make(map[string]interface{}),
		field:    field,
	}

	new.mappings["value"] = value

	return new
}

func (t *Term) Boost(value float64) *Term {
	t.mappings["boost"] = value
	return t
}

func (t *Term) Data() interface{} {
	data := map[string]map[string]interface{}{"term": t.mappings}
	return data
}

func (t *Term) Bytes() []byte {
	bytes, _ := json.Marshal(t.Data())

	return bytes
}

func (t *Term) String() string {
	return string(t.Bytes())
}
