package elastic

import "encoding/json"

type order string
type sortMode string
type missing string

const (
	OrderAsc  order = "asc"
	OrderDesc order = "desc"

	SortModeMin    sortMode = "min"
	SortModeMax    sortMode = "max"
	SortModeSum    sortMode = "sum"
	SortModeAvg    sortMode = "avg"
	SortModeMedian sortMode = "median"

	MissingFirst missing = "_first"
	MissingLast  missing = "_last"
)

type SortField struct {
	mappings map[string]interface{}
	field    string
}

func NewSortField(field string, order ...order) *SortField {
	new := &SortField{
		mappings: make(map[string]interface{}),
		field:    field,
	}

	if order != nil && len(order) > 0 {
		new.Order(order[0])
	}

	return new
}

func (s *SortField) Order(value order) *SortField {
	s.mappings["order"] = value
	return s
}

func (s *SortField) Mode(value sortMode) *SortField {
	s.mappings["mode"] = value
	return s
}

// TODO: needs development
func (s *SortField) Nested(value string) *SortField {
	s.mappings["nested"] = value
	return s
}

func (s *SortField) Missing(value missing) *SortField {
	s.mappings["missing"] = value
	return s
}

func (s *SortField) UnmappedType(value string) *SortField {
	s.mappings["unmapped_type"] = value
	return s
}

// TODO: needs development
func (s *SortField) GgeoDistance(value string) *SortField {
	s.mappings["_geo_distance"] = value
	return s
}

// TODO: needs development
func (s *SortField) Script(value string) *SortField {
	s.mappings["_script"] = value
	return s
}

func (s *SortField) Data() interface{} {
	data := map[string]interface{}{s.field: s.mappings}
	return data
}

func (s *SortField) Bytes() []byte {
	bytes, _ := json.Marshal(s.Data)

	return bytes
}

func (s *SortField) String() string {
	return string(s.Bytes())
}
