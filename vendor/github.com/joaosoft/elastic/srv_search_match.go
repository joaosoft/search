package elastic

import "encoding/json"

type operator string
type zeroTermsQuery string

const (
	OperatorAnd operator = "and"
	OperatorOr  operator = "or"

	ZeroTermsQueryAll  zeroTermsQuery = "all"
	ZeroTermsQueryNone zeroTermsQuery = "none"
)

type Match struct {
	mappings map[string]interface{}
	field    string
}

func NewMatch(field string, query string) *Match {
	new := &Match{
		mappings: make(map[string]interface{}),
		field:    field,
	}

	new.mappings["query"] = query

	return new
}

func (m *Match) Operator(value operator) *Match {
	m.mappings["operator"] = value
	return m
}

func (m *Match) MinimumShouldMatch(value string) *Match {
	m.mappings["minimum_should_match"] = value
	return m
}

func (m *Match) Analyzer(value string) *Match {
	m.mappings["analyzer"] = value
	return m
}

func (m *Match) Lenient(value bool) *Match {
	m.mappings["lenient"] = value
	return m
}

func (m *Match) Fuzziness(value string) *Match {
	m.mappings["fuzziness"] = value
	return m
}

func (m *Match) FuzzyRewrite(value string) *Match {
	m.mappings["fuzzy_rewrite"] = value
	return m
}

func (m *Match) FuzzyTranspositions(value bool) *Match {
	m.mappings["fuzzy_transpositions"] = value
	return m
}

func (m *Match) PrefixLength(value int64) *Match {
	m.mappings["prefix_length"] = value
	return m
}

func (m *Match) MaxExpansions(value int64) *Match {
	m.mappings["max_expansions"] = value
	return m
}

func (m *Match) ZeroTermsQuery(value zeroTermsQuery) *Match {
	m.mappings["zero_terms_query"] = value
	return m
}

func (m *Match) CutoffFrequency(value float64) *Match {
	m.mappings["cutoff_frequency"] = value
	return m
}

func (m *Match) AutoGenerateSynonymsPhraseQuery(value bool) *Match {
	m.mappings["auto_generate_synonyms_phrase_query"] = value
	return m
}

func (m *Match) Data() interface{} {
	data := map[string]map[string]interface{}{"match": {m.field: m.mappings}}
	return data
}

func (m *Match) Bytes() []byte {
	bytes, _ := json.Marshal(m.Data())

	return bytes
}

func (m *Match) String() string {
	return string(m.Bytes())
}
