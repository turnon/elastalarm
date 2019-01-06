package response

import "encoding/json"

type Response struct {
	ScrollID     string          `json:"_scroll_id"`
	Hits         hits            `json:"hits"`
	Aggregations json.RawMessage `json:"aggregations"`
}

type hits struct {
	Total int   `json:"total"`
	Hits  []hit `json:"hits"`
}

type hit struct {
	Index  string          `json:"_index"`
	ID     string          `json:"_id"`
	Source json.RawMessage `json:"_source"`
}

func (resp *Response) Unmarshal(js []byte) {
	if err := json.Unmarshal(js, resp); err != nil {
		panic(err)
	}
}

func (resp *Response) Total() int {
	return resp.Hits.Total
}

func (resp *Response) Aggs() string {
	return string(resp.Aggregations)
}
