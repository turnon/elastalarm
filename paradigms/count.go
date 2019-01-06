package paradigms

import (
	"encoding/json"
	"math/big"

	"github.com/turnon/elastalarm/response"
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

func (c *Count) Found(resp *response.Response) (bool, *string) {
	total := big.NewFloat(float64(resp.Total()))
	match := c.Match.ing(total)
	if !match {
		return match, nil
	}

	detail := "shit"
	return match, &detail
}

func (c *Count) ScopeString() string {
	return string(c.Scope)
}
