package response

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

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

func (resp *Response) Unmarshal(js []byte) error {
	return json.Unmarshal(js, resp)
}

func (resp *Response) Total() int {
	return resp.Hits.Total
}

func (resp *Response) Aggs() string {
	return string(resp.Aggregations)
}

func (resp *Response) FlattenAggs() string {
	var sb strings.Builder

	resp.FlatEach(func(arr []interface{}, count int) {
		sb.WriteString(fmt.Sprint(arr))
		sb.WriteString(" ")
		sb.WriteString(strconv.Itoa(count))
		sb.WriteString("\n")
	})

	return sb.String()
}

func (resp *Response) FlatEach(f func([]interface{}, int)) {
	var b bucket
	b.unmarshal(resp.Aggregations)
	b.flatten(f)
}
