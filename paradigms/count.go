package paradigms

import (
	"fmt"
	"math/big"

	"github.com/turnon/elastalarm/response"
)

type Count struct {
	EsDsl
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
				{{ .Paradigm.QueryString }}
			]
		}
	},
	"size": 0,
	"aggs": {{ .Paradigm.AggsString }}
}
`

func (c *Count) Template() string {
	return countTemplate
}

func (c *Count) SupportStep() bool {
	return true
}

func (c *Count) Found(resp *response.Response) (bool, *response.Result) {
	total := big.NewFloat(float64(resp.Total()))
	match, desc := c.match(total)
	if !match {
		return match, nil
	}

	result := &response.Result{}
	resp.FlatEach(func(arr []interface{}, count int) {
		result.SetDetail(arr, count, nil)
	})
	abstract := fmt.Sprintf("total %d %s", resp.Total(), desc)
	result.Abstract = abstract

	return match, result
}

func (c *Count) FoundOnAggs(resp *response.Response) (bool, *response.Result) {
	var (
		anyMatch bool
		anyDesc  string
	)

	result := &response.Result{}

	resp.FlatEach(func(arr []interface{}, count int) {
		total := big.NewFloat(float64(count))
		if match, desc := c.match(total); match {
			anyMatch = match
			anyDesc = desc
			result.SetDetail(arr, count, nil)
		}
	})

	if !anyMatch {
		return false, nil
	}

	abstract := fmt.Sprintf("something %s", anyDesc)
	result.Abstract = abstract

	return anyMatch, result
}
