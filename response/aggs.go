package response

import (
	"encoding/json"
	"strconv"
)

type singleAggName struct {
	prefix, name string
}

type singleAgg struct {
	DocCountErrorUpperBound int             `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int             `json:"sum_other_doc_count"`
	RawBuckets              json.RawMessage `json:"buckets"`
	_buckets                []*bucket
}

func (a *singleAgg) buckets() []*bucket {
	if a._buckets == nil {
		var (
			bs         []*bucket
			rawBuckets []json.RawMessage
		)
		json.Unmarshal(a.RawBuckets, &rawBuckets)
		for _, rawBucket := range rawBuckets {
			b := &bucket{}
			b.unmarshal(rawBucket)
			bs = append(bs, b)
		}
		a._buckets = bs
	}

	return a._buckets
}

const wrongKey = "wrong_key"

type bucket struct {
	key             string
	docCount        int
	nameRawAggPairs map[string](*json.RawMessage)
	nameAggPairs    map[singleAggName](*singleAgg)
}

func (b *bucket) unmarshal(bytes []byte) {
	var nameRawAggPairs map[string](*json.RawMessage)
	json.Unmarshal(bytes, &nameRawAggPairs)

	for key, rawMsg := range nameRawAggPairs {
		if key == "key" {
			b.key = string(*rawMsg)
		} else if key == "doc_count" {
			count, _ := strconv.Atoi(string(*rawMsg))
			b.docCount = count
		} else {
			if b.nameRawAggPairs == nil {
				b.nameRawAggPairs = make(map[string](*json.RawMessage))
			}
			b.nameRawAggPairs[key] = rawMsg
		}
	}
}

func (b *bucket) aggs() *map[singleAggName]*singleAgg {
	if b.nameAggPairs == nil {
		b.cacheAggs()
	}
	return &b.nameAggPairs
}

func (b *bucket) cacheAggs() {
	aggs := make(map[singleAggName]*singleAgg)

	for key, rawMsgAddr := range b.nameRawAggPairs {
		var tryUnmarshal map[string]*json.RawMessage
		json.Unmarshal(*rawMsgAddr, &tryUnmarshal)

		if tryUnmarshal["doc_count"] == nil {
			name := singleAggName{name: key}
			agg := &singleAgg{}
			json.Unmarshal(*rawMsgAddr, agg)
			aggs[name] = agg
		} else {
			var filter bucket
			filter.unmarshal(*rawMsgAddr)
			for name, agg := range *filter.aggs() {
				name.prefix = key
				aggs[name] = agg
			}
		}
	}

	b.nameAggPairs = aggs
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

	for name, a := range aggs {
		for _, b := range a.buckets() {
			var moreKeys = keys
			if name.prefix != "" {
				moreKeys = append(moreKeys, name.prefix)
			}
			moreKeys = append(moreKeys, b.key)
			b._flatten(moreKeys, b.docCount, f)
		}
	}
}
