package elastic

import "encoding/json"

type Boosting struct {
	mappings      map[string]interface{}
}

func NewBoosting() *Boosting {
	new := &Boosting{
		mappings: make(map[string]interface{}),
	}

	return new
}

func (b *Boosting) Positive(value ...Query) *Boosting {
	if len(value) > 0 {
		b.mappings["positive"] = append(b.mappings["positive"].([]Query), value...)
	}
	return b
}

func (b *Boosting) Negative(value ...Query) *Boosting {
	if len(value) > 0 {
		b.mappings["negative"] = append(b.mappings["negative"].([]Query), value...)
	}
	return b
}

func (b *Boosting) NegativeBoost(value float64) *Boosting {
	b.mappings["negative_boost"] = value
	return b
}

func (b *Boosting) Data() interface{} {
	data := make(map[string]interface{})

	if len(b.mappings) > 0 {
		data["boosting"] = b.mappings
	}

	return data
}

func (b *Boosting) Bytes() []byte {
	bytes, _ := json.Marshal(b.Data())

	return bytes
}

func (b *Boosting) String() string {
	return string(b.Bytes())
}
