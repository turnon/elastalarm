package paradigms

import (
	"encoding/json"
	"math/big"

	"github.com/turnon/elastalarm/response"
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
			"aggs": {{ .DetailString }}
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

	detail := resp.FlattenAggs()
	return match, &detail
}

func (p *Percentage) PartString() string {
	return string(p.Part)
}

func (p *Percentage) WholeString() string {
	return string(p.Whole)
}
