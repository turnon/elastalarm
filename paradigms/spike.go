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

func (s *Spike) Found(resp *response.Response) (bool, *string) {
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

	detail := fmt.Sprintf("%d / %d = %s %s. actual past doc_ount is %d \n\n%s",
		aggs.Recent.DocCount, pastCount, times.String(), desc, aggs.Past.DocCount, resp.FlattenAggs())
	return match, &detail
}

func (s *Spike) FoundOnAggs(resp *response.Response) (bool, *string) {
	var (
		anyMatch bool
		anyDesc  string
	)

	formator := response.GetFormator("")()
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
		// fmt.Println(past, "-------")
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
			formator.SetDetail(keys, count)
		}
	})

	// calculate 0/past if past remain
	var recentNotFound int
	if s.Ref == 0 {
		recentNotFound = 0
	} else {
		recentNotFound = s.Ref
	}

	recentNotFoundFloat := big.NewFloat(float64(recentNotFound))
	if match, desc := s.match(recentNotFoundFloat); match {
		anyMatch = match
		anyDesc = desc
		for key, count := range past {
			keys := pastRawKeys[key]
			formator.SetDetail(keys, count)
		}
	}

	if !anyMatch {
		return false, nil
	}

	abstract := fmt.Sprintf("something %s", anyDesc)
	formator.SetAbstract(abstract)
	detail := formator.String()
	return anyMatch, &detail
}

func calcTimes(past, recent int) *big.Float {
	pastF := big.NewFloat(float64(past))
	recentF := big.NewFloat(float64(recent))

	var times big.Float
	times.Quo(pastF, recentF)

	return &times
}
