package paradigms

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/turnon/elastalarm/response"
)

type Spike struct {
	Scope json.RawMessage `json:"scope"`
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
						"@timestamp": {
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
					"@timestamp": {
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
					"@timestamp": {
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

func (s *Spike) Found(resp *response.Response) (bool, *string) {
	aggs := &spikeAggs{}
	json.Unmarshal(resp.Aggregations, aggs)
	past := big.NewFloat(float64(aggs.Past.DocCount))
	recent := big.NewFloat(float64(aggs.Recent.DocCount))

	var times big.Float
	times.Quo(recent, past)
	match, desc := s.match(&times)

	if !match {
		return match, nil
	}

	detail := fmt.Sprintf("%d / %d = %s %s\n\n%s",
		aggs.Recent.DocCount, aggs.Past.DocCount, times.String(), desc, resp.FlattenAggs())
	return match, &detail
}

func (s *Spike) ScopeString() string {
	return string(s.Scope)
}