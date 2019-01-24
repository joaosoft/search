package search

type direction string

const (
	orderAsc  direction = "asc"
	orderDesc direction = "desc"
)

type order struct {
	column    string
	direction direction
}

type orders []*order
