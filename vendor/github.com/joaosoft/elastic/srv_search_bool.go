package elastic

import "encoding/json"

type Bool struct {
	mappings  map[string]interface{}
}

func NewBool() *Bool {
	new := &Bool{
		mappings: make(map[string]interface{}),
	}

	return new
}

func (b *Bool) Must(value ...Query) *Bool {
	if len(value) > 0 {
		b.mappings["must"] = append(b.mappings["must"].([]Query), value...)
	}
	return b
}

func (b *Bool) MustNot(value ...Query) *Bool {
	if len(value) > 0 {
		b.mappings["must_not"] = append(b.mappings["must_not"].([]Query), value...)
	}
	return b
}

func (b *Bool) Filter(value ...Query) *Bool {
	if len(value) > 0 {
		b.mappings["filter"] = append(b.mappings["filter"].([]Query), value...)
	}
	return b
}

func (b *Bool) Should(value ...Query) *Bool {
	if len(value) > 0 {
		b.mappings["should"] = append(b.mappings["should"].([]Query), value...)
	}
	return b
}

func (b *Bool) Data() interface{} {
	data := make(map[string]interface{})

	if len(b.mappings) > 0 {
		data["bool"] = b.mappings
	}

	return data
}

func (b *Bool) Bytes() []byte {
	bytes, _ := json.Marshal(b.Data())

	return bytes
}

func (b *Bool) String() string {
	return string(b.Bytes())
}
