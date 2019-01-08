package paradigms

import (
	"encoding/json"
	"math/big"

	"github.com/turnon/elastalarm/response"
)

type Spike struct {
	Scope json.RawMessage `json:"scope"`
	Match `json:"match"`
}

type spikeAggs struct {
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
							"gt": "now-{{ .Interval }}-{{ .Interval }}"
						}
					}
				},
				{{ .Paradigm.ScopeString }}
			]
		}
	},
	"size": 0,
	"aggs": {
		"recent": {
			"filter": {
				"range": {
					"@timestamp": {
						"gt": "now-{{ .Interval }}"
					}
				}
			},
			"aggs": {{ .DetailString }}
		}
	}
}
`

func (s *Spike) Template() string {
	return percentageTemplate
}

func (s *Spike) Found(resp *response.Response) (bool, *string) {
	pastAndRecent := resp.Total()
	aggs := &spikeAggs{}
	json.Unmarshal(resp.Aggregations, aggs)
	past := big.NewFloat(float64(pastAndRecent - aggs.Recent.DocCount))
	recent := big.NewFloat(float64(aggs.Recent.DocCount))

	var times big.Float
	times.Quo(recent, past)
	match := s.Match.ing(&times)

	if !match {
		return match, nil
	}

	detail := resp.FlattenAggs()
	return match, &detail
}

func (s *Spike) ScopeString() string {
	return string(s.Scope)
}
