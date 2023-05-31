package elastic

import "encoding/json"

type QueryString struct {
	mappings map[string]interface{}
}

func NewQueryString(query interface{}) *QueryString {
	new := &QueryString{
		mappings: map[string]interface{}{"query": query},
	}

	return new
}

func (q *QueryString) DefaultField(value float64) *QueryString {
	q.mappings["default_field"] = value
	return q
}

func (q *QueryString) DefaultOperator(value operator) *QueryString {
	q.mappings["default_operator"] = value
	return q
}

func (q *QueryString) Analyser(value string) *QueryString {
	q.mappings["analyser"] = value
	return q
}

func (q *QueryString) QuoteAnalyzer(value string) *QueryString {
	q.mappings["quote_analyzer"] = value
	return q
}

func (q *QueryString) AllowLeadingWildcard(value bool) *QueryString {
	q.mappings["allow_leading_wildcard"] = value
	return q
}

func (q *QueryString) EnablePositionIncrements(value bool) *QueryString {
	q.mappings["enable_position_increments"] = value
	return q
}

func (q *QueryString) FuzzyMaxExpansions(value int64) *QueryString {
	q.mappings["fuzzy_max_expansions"] = value
	return q
}

func (q *QueryString) Fuzziness(value string) *QueryString {
	q.mappings["fuzziness"] = value
	return q
}

func (q *QueryString) FuzzyPrefixLength(value int64) *QueryString {
	q.mappings["fuzzy_prefix_length"] = value
	return q
}

func (q *QueryString) FuzzyTranspositions(value bool) *QueryString {
	q.mappings["fuzzy_transpositions"] = value
	return q
}

func (q *QueryString) PhraseSlop(value int64) *QueryString {
	q.mappings["phrase_slop"] = value
	return q
}

func (q *QueryString) Boost(value int64) *QueryString {
	q.mappings["boost"] = value
	return q
}

func (q *QueryString) AutoGeneratePhraseQueries(value bool) *QueryString {
	q.mappings["auto_generate_phrase_queries"] = value
	return q
}

func (q *QueryString) AnalyzeWildcard(value bool) *QueryString {
	q.mappings["analyze_wildcard"] = value
	return q
}

func (q *QueryString) MaxDeterminizedStates(value bool) *QueryString {
	q.mappings["max_determinized_states"] = value
	return q
}

func (q *QueryString) MinimumShouldMatch(value string) *QueryString {
	q.mappings["minimum_should_match"] = value
	return q
}

func (q *QueryString) Lenient(value bool) *QueryString {
	q.mappings["lenient"] = value
	return q
}

func (q *QueryString) TimeZone(value string) *QueryString {
	q.mappings["time_zone"] = value
	return q
}

func (q *QueryString) QuoteFieldSuffix(value string) *QueryString {
	q.mappings["quote_field_suffix"] = value
	return q
}

func (q *QueryString) AutoGenerateSynonymsPhraseQuery(value bool) *QueryString {
	q.mappings["auto_generate_synonyms_phrase_query"] = value
	return q
}

func (q *QueryString) AllFields(value string) *QueryString {
	q.mappings["all_fields"] = value
	return q
}

func (q *QueryString) TieBreaker(value string) *QueryString {
	q.mappings["tie_breaker"] = value
	return q
}

func (q *QueryString) Fields(values ...string) *QueryString {
	if _, ok := q.mappings["fields"]; !ok {
		q.mappings["fields"] = make([]string, 0)
	}

	q.mappings["fields"] = append(q.mappings["fields"].([]string), values...)
	return q
}

func (q *QueryString) Data() interface{} {
	data := map[string]interface{}{"query_string": q.mappings}
	return data
}

func (q *QueryString) Bytes() []byte {
	bytes, _ := json.Marshal(q.Data())

	return bytes
}

func (q *QueryString) String() string {
	return string(q.Bytes())
}
