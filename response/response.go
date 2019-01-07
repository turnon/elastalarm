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

type singleAgg struct {
	DocCountErrorUpperBound int      `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int      `json:"sum_other_doc_count"`
	Buckets                 []bucket `json:"buckets"`
}

type bucket map[string](*json.RawMessage)

func (b *bucket) key() string {
	if rawMsgAddr := (*b)["key"]; rawMsgAddr != nil {
		return string(*rawMsgAddr)
	}
	return ""
}

func (b *bucket) docCount() int {
	if rawMsgAddr := (*b)["doc_count"]; rawMsgAddr != nil {
		i, _ := strconv.Atoi(string(*rawMsgAddr))
		return i
	}
	return 0
}

func (b *bucket) aggs() *map[string]*singleAgg {
	aggs := make(map[string]*singleAgg)

	for key, rawMsgAddr := range *b {
		if key == "key" || key == "doc_count" {
			continue
		}

		agg := &singleAgg{}
		json.Unmarshal(*rawMsgAddr, agg)
		aggs[key] = agg
	}

	return &aggs
}

func (b *bucket) flatten(f func([]interface{}, int)) {
	var keys []interface{}
	b._flatten(keys, 0, f)
}

func (b *bucket) _flatten(keys []interface{}, count int, f func([]interface{}, int)) {
	aggs := *b.aggs()
	if len(aggs) == 0 {
		f(keys, count)
		return
	}

	for _, a := range aggs {
		for _, b := range a.Buckets {
			moreKeys := append(keys, b.key())
			b._flatten(moreKeys, b.docCount(), f)
		}
	}
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

func (resp *Response) FlattenAggs() string {
	var (
		b  bucket
		sb strings.Builder
	)
	json.Unmarshal(resp.Aggregations, &b)

	b.flatten(func(arr []interface{}, count int) {
		sb.WriteString(fmt.Sprint(arr))
		sb.WriteString(" ")
		sb.WriteString(strconv.Itoa(count))
		sb.WriteString("\n")
	})

	return sb.String()
}
