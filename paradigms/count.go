package paradigms

import (
	"encoding/json"
	"fmt"
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
							"gt": "{{ .NowString }}-{{ .Interval }}"
						}
					}
				},
				{{ .Paradigm.ScopeString }}
			]
		}
	},
	"size": 0,
	"aggs": {{ .DetailString }}
}
`

func (c *Count) Template() string {
	return countTemplate
}

func (c *Count) Found(resp *response.Response) (bool, *string) {
	total := big.NewFloat(float64(resp.Total()))
	match, desc := c.match(total)
	if !match {
		return match, nil
	}

	detail := fmt.Sprintf("total %d %s\n\n%s", resp.Total(), desc, resp.FlattenAggs())
	return match, &detail
}

func (c *Count) ScopeString() string {
	return string(c.Scope)
}
