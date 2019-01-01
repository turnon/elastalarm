package paradigms

import (
	"encoding/json"
	"math/big"
	"regexp"

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
				{
					{{ .Paradigm.WholeString }}
				}
			]
	  }
	},
	"size": 0,
	"aggs": {
	  "part": {
			"filter": {
				{{ .Paradigm.PartString }}
			},
			"aggs": {
				{{ .Paradigm.DetailString }}
			}
	  }
	}
}
`

var (
	re      = regexp.MustCompile("(?s)\\{(.*)\\}")
	hundred = big.NewFloat(100)
)

func (p *Percentage) Template() string {
	return percentageTemplate
}

func (p *Percentage) Found(resp *response.Response) bool {
	total := resp.Total()
	if total == 0 {
		return false
	}
	whole := big.NewFloat(float64(total))

	aggs := &percentAggs{}
	json.Unmarshal(resp.Aggregations, aggs)
	part := big.NewFloat(float64(aggs.Part.DocCount))

	var quo, percent big.Float
	quo.Quo(part, whole)
	percent.Mul(&quo, hundred)
	return p.Match.ing(&percent)
}

func stringify(json *json.RawMessage) string {
	return re.ReplaceAllString(string(*json), `$1`)
}

func (p *Percentage) PartString() string {
	return stringify(&p.Part)
}

func (p *Percentage) WholeString() string {
	return stringify(&p.Whole)
}

func (p *Percentage) DetailString() string {
	return stringify(&p.Detail)
}
