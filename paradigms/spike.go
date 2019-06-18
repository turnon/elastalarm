package paradigms

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/turnon/elastalarm/response"
)

type Spike struct {
	Scope json.RawMessage `json:"scope"`
	Ref   int             `json:"reference"`
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
				{{ .Paradigm.ScopeString }}
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
			"aggs": {{ .DetailString }}
		},
		"recent": {
			"filter": {
				"range": {
					"{{ .TimeField }}": {
						"gt": "{{ .NowString }}-{{ .Interval }}"
					}
				}
			},
			"aggs": {{ .DetailString }}
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

	past := big.NewFloat(float64(pastCount))
	recent := big.NewFloat(float64(aggs.Recent.DocCount))

	var times big.Float
	times.Quo(recent, past)
	match, desc := s.match(&times)

	if !match {
		return match, nil
	}

	detail := fmt.Sprintf("%d / %d = %s %s. actual past doc_ount is %d \n\n%s",
		aggs.Recent.DocCount, pastCount, times.String(), desc, aggs.Past.DocCount, resp.FlattenAggs())
	return match, &detail
}

func (s *Spike) FoundOnDetail(resp *response.Response) (bool, *string) {
	detail := ""
	return false, &detail
}

func (s *Spike) ScopeString() string {
	return string(s.Scope)
}
