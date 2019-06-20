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
						"{{ .TimeField }}": {
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

	formator := response.GetFormator("")()
	resp.FlattenDetail(formator)
	abstract := fmt.Sprintf("total %d %s", resp.Total(), desc)
	formator.SetAbstract(abstract)

	detail := formator.String()
	return match, &detail
}

func (c *Count) FoundOnDetail(resp *response.Response) (bool, *string) {
	var (
		anyMatch bool
		anyDesc  string
	)

	formator := response.GetFormator("")()

	resp.FlatEach(func(arr []interface{}, count int) {
		total := big.NewFloat(float64(count))
		if match, desc := c.match(total); match {
			anyMatch = match
			anyDesc = desc
			formator.SetDetail(arr, count)
		}
	})

	if !anyMatch {
		return false, nil
	}

	abstract := fmt.Sprintf("something %s", anyDesc)
	formator.SetAbstract(abstract)
	detail := formator.String()
	return anyMatch, &detail
}

func (c *Count) ScopeString() string {
	return string(c.Scope)
}
