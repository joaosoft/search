package elastic

import "encoding/json"

type Sort struct {
	mappings []*SortField
}

func NewSort(fields ...*SortField) *Sort {
	new := &Sort{
		mappings: fields,
	}

	return new
}

func (s *Sort) Data() interface{} {
	data := make(map[string]interface{})

	if len(s.mappings) > 0 {
		data["sort"] =  s.mappings
	}

	return data
}

func (s *Sort) Bytes() []byte {
	bytes, _ := json.Marshal(s.Data())

	return bytes
}

func (s *Sort) String() string {
	return string(s.Bytes())
}
