package paradigms

import (
	"encoding/json"
	"math/big"

	"bitbucket.org/xcrossing/elastic_alarm/response"
)

type Count struct {
	Scope json.RawMessage `json:"scope"`
	Match `json:"match"`
}

const countTemplate = `
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
				{{ .Paradigm.ScopeString }}
			]
		}
	},
	"size": 0
}
`

func (c *Count) Template() string {
	return countTemplate
}

func (c *Count) Found(resp *response.Response) bool {
	total := big.NewFloat(float64(resp.Total()))
	return c.Match.ing(total)
}

func (c *Count) ScopeString() string {
	return string(c.Scope)
}
