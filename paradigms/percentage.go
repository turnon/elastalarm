package paradigms

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/turnon/elastalarm/response"
)

type Percentage struct {
	EsDsl
	PartialQuery json.RawMessage `json:"partial_query"`
	Match        `json:"match"`
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
	"aggs": {
		"part": {
			"filter": {{ .Paradigm.PartialQueryString }},
			"aggs": {{ .Paradigm.AggsString }}
		}
	}
}
`

const percentageOnAggsTemplate = `
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

var hundred = big.NewFloat(100)

func (p *Percentage) Template() string {
	if p.OnAggs() {
		return percentageOnAggsTemplate
	}
	return percentageTemplate
}

func (p *Percentage) Found(resp *response.Response) (bool, *response.Result) {
	total := resp.Total()
	if total == 0 {
		return false, nil
	}

	aggs := &percentAggs{}
	json.Unmarshal(resp.Aggregations, aggs)
	part := aggs.Part.DocCount

	percent := calcPercent(part, total)
	match, desc := p.match(percent)

	if !match {
		return match, nil
	}

	result := &response.Result{}
	resp.FlatEach(func(arr []interface{}, count int) {
		result.SetDetail(arr, count, nil)
	})
	abstract := fmt.Sprintf("%d / %d = %s%% %s", part, total, percent.String(), desc)
	result.SetAbstract(abstract)

	return match, result
}

func (p *Percentage) FoundOnAggs(resp *response.Response) (bool, *response.Result) {
	total := resp.Total()
	if total == 0 {
		return false, nil
	}

	var (
		anyMatch bool
		anyDesc  string
	)

	result := &response.Result{}

	resp.FlatEach(func(arr []interface{}, part int) {
		percent := calcPercent(part, total)
		if match, desc := p.match(percent); match {
			anyMatch = match
			anyDesc = desc
			result.SetDetail(arr, part, percent)
		}
	})

	if !anyMatch {
		return false, nil
	}

	abstract := fmt.Sprintf("something %s", anyDesc)
	result.SetAbstract(abstract)
	return anyMatch, result
}

func calcPercent(numerator, denominator int) *big.Float {
	n := big.NewFloat(float64(numerator))
	d := big.NewFloat(float64(denominator))

	var quo, percent big.Float
	quo.Quo(n, d)
	percent.Mul(&quo, hundred)

	return &percent
}

func (p *Percentage) PartialQueryString() string {
	return string(p.PartialQuery)
}
