package paradigms

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/turnon/elastalarm/response"
)

type Spike struct {
	EsDsl
	Ref   int `json:"reference"`
	Match `json:"match"`
}

type spikeAggs struct {
	Past struct {
		DocCount int `json:"doc_count"`
	} `json:"past"`
	Recent struct {
		DocCount int `json:"doc_count"`
	} `json:"recent"`
}

const spikeTemplate = `
{
	"query": {
		"bool": {
			"must": [
				{
					"range": {
						"{{ .TimeField }}": {
							"gt": "{{ .NowString }}-{{ .Interval }}-{{ .Interval }}"
						}
					}
				},
				{{ .Paradigm.QueryString }}
			]
		}
	},
	"size": 0,
	"aggs": {
		"past": {
			"filter": {
				"range": {
					"{{ .TimeField }}": {
						"gt": "{{ .NowString }}-{{ .Interval }}-{{ .Interval }}",
						"lte": "{{ .NowString }}-{{ .Interval }}"
					}
				}
			},
			"aggs": {{ .Paradigm.AggsString }}
		},
		"recent": {
			"filter": {
				"range": {
					"{{ .TimeField }}": {
						"gt": "{{ .NowString }}-{{ .Interval }}"
					}
				}
			},
			"aggs": {{ .Paradigm.AggsString }}
		}
	}
}
`

func (s *Spike) Template() string {
	return spikeTemplate
}

func (s *Spike) pastCount(aggs *spikeAggs) int {
	count := aggs.Past.DocCount
	if count != 0 {
		return count
	}

	if s.Ref != 0 {
		return s.Ref
	}

	return 0
}

func (s *Spike) Found(resp *response.Response) (bool, *response.Result) {
	aggs := &spikeAggs{}
	json.Unmarshal(resp.Aggregations, aggs)

	pastCount := s.pastCount(aggs)
	if pastCount == 0 {
		return false, nil
	}

	times := calcTimes(pastCount, aggs.Recent.DocCount)
	match, desc := s.match(times)

	if !match {
		return match, nil
	}

	result := &response.Result{}
	resp.FlatEach(func(arr []interface{}, count int) {
		result.SetDetail(arr, count, nil)
	})
	abstract := fmt.Sprintf("%d / %d = %s %s. actual past doc_ount is %d",
		aggs.Recent.DocCount, pastCount, times.String(), desc, aggs.Recent.DocCount)
	result.Abstract = abstract

	return match, result
}

func (s *Spike) FoundOnAggs(resp *response.Response) (bool, *response.Result) {
	var (
		anyMatch bool
		anyDesc  string
	)

	result := &response.Result{}
	past := make(map[string]int)
	pastRawKeys := make(map[string][]interface{})

	resp.FlatEach(func(arr []interface{}, count int) {
		from := fmt.Sprint(arr[0])
		keys := arr[1:]
		keystr := fmt.Sprint(keys)

		// cache past count
		if from == "past" {
			past[keystr] = count
			pastRawKeys[keystr] = keys
			return
		}

		// calculate recent/past and remove key in both recent and past
		pastCount := past[keystr]
		delete(past, keystr)
		if pastCount == 0 {
			if s.Ref == 0 {
				return
			}
			pastCount = s.Ref
		}

		times := calcTimes(pastCount, count)
		if match, desc := s.match(times); match {
			anyMatch = match
			anyDesc = desc
			result.SetDetail(keys, count, times)
		}
	})

	// calculate 0/past if past remain
	if s.Ref == 0 {
		recentNotFound := big.NewFloat(float64(0))
		if match, desc := s.match(recentNotFound); match {
			anyMatch = match
			anyDesc = desc
			for key, _ := range past {
				keys := pastRawKeys[key]
				result.SetDetail(keys, 0, recentNotFound)
			}
		}
	} else {
		for key, count := range past {
			times := calcTimes(count, s.Ref)
			if match, desc := s.match(times); match {
				anyMatch = match
				anyDesc = desc
				keys := pastRawKeys[key]
				result.SetDetail(keys, s.Ref, times)
			}
		}
	}

	if !anyMatch {
		return false, nil
	}

	abstract := fmt.Sprintf("something %s", anyDesc)
	result.Abstract = abstract

	return anyMatch, result
}

func calcTimes(past, recent int) *big.Float {
	pastF := big.NewFloat(float64(past))
	recentF := big.NewFloat(float64(recent))

	var times big.Float
	times.Quo(recentF, pastF)

	return &times
}
