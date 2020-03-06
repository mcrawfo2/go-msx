package lucene

import "encoding/json"

type SearchBuilder struct {
	Filter  []Condition `json:"filter,omitempty"`
	Query   []Condition `json:"query,omitempty"`
	Sort    []SortField `json:"sort,omitempty"`
	Refresh *bool       `json:"refresh,omitempty"`
}

func (s *SearchBuilder) WithRefresh(refresh bool) *SearchBuilder {
	s.Refresh = &refresh
	return s
}

func (s *SearchBuilder) WithFilter(filter ...Condition) *SearchBuilder {
	s.Filter = append(s.Filter, filter...)
	return s
}

func (s *SearchBuilder) WithQuery(query ...Condition) *SearchBuilder {
	s.Query = append(s.Query, query...)
	return s
}

func (s *SearchBuilder) WithSort(sort ...SortField) *SearchBuilder {
	s.Sort = append(s.Sort, sort...)
	return s
}

func (s *SearchBuilder) Build() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}

type Condition interface {
	ConditionString() string
}

type BooleanCondition struct {
	Type   string      `json:"type"`
	Boost  *float32    `json:"boost,omitempty"`
	Must   []Condition `json:"must,omitempty"`
	Should []Condition `json:"should,omitempty"`
	Not    []Condition `json:"not,omitempty"`
}

func (b *BooleanCondition) ConditionString() string {
	bytes, _ := json.Marshal(b)
	return string(bytes)
}

func (b *BooleanCondition) WithBoost(boost float32) *BooleanCondition {
	b.Boost = &boost
	return b
}

func (b *BooleanCondition) WithMust(condition ...Condition) *BooleanCondition {
	b.Must = append(b.Must, condition...)
	return b
}

func (b *BooleanCondition) WithShould(condition ...Condition) *BooleanCondition {
	b.Should = append(b.Should, condition...)
	return b
}

func (b *BooleanCondition) WithNot(condition ...Condition) *BooleanCondition {
	b.Not = append(b.Not, condition...)
	return b
}

func NewBoolean() *BooleanCondition {
	return &BooleanCondition{
		Type: "boolean",
	}
}

func Must(conditions ...Condition) *BooleanCondition {
	return NewBoolean().WithMust(conditions...)
}

func Should(conditions ...Condition) *BooleanCondition {
	return NewBoolean().WithShould(conditions...)
}

func Not(conditions ...Condition) *BooleanCondition {
	return NewBoolean().WithNot(conditions...)
}

type MatchCondition struct {
	Type      string      `json:"type"`
	Boost     float32     `json:"boost,omitempty"`
	Field     string      `json:"field"`
	Value     interface{} `json:"value"`
	DocValues *bool       `json:"doc_values,omitempty"`
}

func (m *MatchCondition) WithBoost(boost float32) *MatchCondition {
	m.Boost = boost
	return m
}

func (m *MatchCondition) WithField(field string) *MatchCondition {
	m.Field = field
	return m
}

func (m *MatchCondition) WithValue(value interface{}) *MatchCondition {
	m.Value = value
	return m
}

func (m *MatchCondition) WithDocValues(docValues bool) *MatchCondition {
	m.DocValues = &docValues
	return m
}

func (m *MatchCondition) ConditionString() string {
	bytes, _ := json.Marshal(m)
	return string(bytes)
}

func NewMatchCondition() *MatchCondition {
	return &MatchCondition{
		Type: "match",
	}
}

func Match(field string, value interface{}) *MatchCondition {
	return NewMatchCondition().WithField(field).WithValue(value)
}

type RangeCondition struct {
	Type         string      `json:"type"`
	Boost        float32     `json:"boost,omitempty"`
	Field        string      `json:"field"`
	Lower        interface{} `json:"lower,omitempty"`
	Upper        interface{} `json:"upper,omitempty"`
	IncludeLower *bool       `json:"include_lower,omitempty"`
	IncludeUpper *bool       `json:"include_upper,omitempty"`
	DocValues    *bool       `json:"doc_values,omitempty"`
}

func (r *RangeCondition) WithField(field string) *RangeCondition {
	r.Field = field
	return r
}

func (r *RangeCondition) WithLower(lower interface{}) *RangeCondition {
	r.Lower = lower
	return r
}

func (r *RangeCondition) WithUpper(upper interface{}) *RangeCondition {
	r.Upper = upper
	return r
}

func (r *RangeCondition) WithIncludeLower(includeLower bool) *RangeCondition {
	r.IncludeLower = &includeLower
	return r
}

func (r *RangeCondition) WithIncludeUpper(includeUpper bool) *RangeCondition {
	r.IncludeUpper = &includeUpper
	return r
}

func (r *RangeCondition) WithDocValues(docValues bool) *RangeCondition {
	r.DocValues = &docValues
	return r
}

func (r *RangeCondition) ConditionString() string {
	bytes, _ := json.Marshal(r)
	return string(bytes)
}

func NewRangeCondition() *RangeCondition {
	return &RangeCondition{
		Type: "range",
	}
}

func Range(field string) *RangeCondition {
	return NewRangeCondition().WithField(field)
}

type SortField interface {
	SortString() string
}

type SimpleSortField struct {
	Type    string `json:"type"`
	Reverse *bool  `json:"reverse,omitempty"`
	Field   string `json:"field"`
}

func (s *SimpleSortField) WithReverse(reverse bool) *SimpleSortField {
	s.Reverse = &reverse
	return s
}

func (s *SimpleSortField) WithField(field string) *SimpleSortField {
	s.Field = field
	return s
}

func (s *SimpleSortField) SortString() string {
	bytes, _ := json.Marshal(&s)
	return string(bytes)
}

func NewSimpleSortField() *SimpleSortField {
	return &SimpleSortField{Type: "simple"}
}

func Field(field string) *SimpleSortField {
	return NewSimpleSortField().WithField(field)
}
