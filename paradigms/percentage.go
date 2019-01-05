package paradigms

import (
	"encoding/json"
	"math/big"

	"bitbucket.org/xcrossing/elastic_alarm/response"
)

type Percentage struct {
	Part   json.RawMessage `json:"part"`
	Whole  json.RawMessage `json:"whole"`
	Match  `json:"match"`
	Detail json.RawMessage `json:"detail"`
}

type percentAggs struct {
	Part struct {
		DocCount int `json:"doc_count"`
	} `json:"part"`
}

const percentageTemplate = `
{
	"query": {
		"bool": {
			"must": [
				{
					"range": {
						"@timestamp": {
							"gt": "now-{{ .Interval }}"
						}
					}
				},
				{{ .Paradigm.WholeString }}
			]
		}
	},
	"size": 0,
	"aggs": {
		"part": {
			"filter": {{ .Paradigm.PartString }},
			"aggs": {
				{{ .Paradigm.DetailString }}
			}
		}
	}
}
`

var hundred = big.NewFloat(100)

func (p *Percentage) Template() string {
	return percentageTemplate
}

func (p *Percentage) Found(resp *response.Response) (bool, *string) {
	total := resp.Total()
	if total == 0 {
		return false, nil
	}
	whole := big.NewFloat(float64(total))

	aggs := &percentAggs{}
	json.Unmarshal(resp.Aggregations, aggs)
	part := big.NewFloat(float64(aggs.Part.DocCount))

	var quo, percent big.Float
	quo.Quo(part, whole)
	percent.Mul(&quo, hundred)
	match := p.Match.ing(&percent)

	if !match {
		return match, nil
	}

	detail := "wtf"
	return match, &detail
}

func (p *Percentage) PartString() string {
	return string(p.Part)
}

func (p *Percentage) WholeString() string {
	return string(p.Whole)
}

func (p *Percentage) DetailString() string {
	return string(p.Detail)
}
