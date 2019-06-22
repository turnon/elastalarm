package paradigms

import "encoding/json"

type EsDsl struct {
	Query json.RawMessage `json:"query"`
	Aggs  json.RawMessage `json:"aggs"`
}

func (dsl *EsDsl) QueryString() string {
	if str := string(dsl.Query); str != "" {
		return str
	}
	return `{"match_all": {}}`
}

func (dsl *EsDsl) AggsString() string {
	if str := string(dsl.Aggs); str != "" {
		return str
	}
	return `{}`
}
