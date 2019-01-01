package paradigms

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"text/template"

	"bitbucket.org/xcrossing/elastic_alarm/response"
)

type Percentage struct {
	Part   json.RawMessage `json:"part"`
	Whole  json.RawMessage `json:"whole"`
	Gt     float64         `json:"gt"`
	Lt     float64         `json:"lt"`
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
							"gt": "now-30d"
						}
					}
				},
				{
					{{ .WholeString }}
				}
			]
	  }
	},
	"size": 0,
	"aggs": {
	  "part": {
			"filter": {
				{{ .PartString }}
			},
			"aggs": {
				{{ .DetailString }}
			}
	  }
	}
}
`

var re = regexp.MustCompile("(?s)\\{(.*)\\}")

func (p *Percentage) ReqBody() io.Reader {
	t := template.New("a")
	t.Parse(percentageTemplate)
	s := &strings.Builder{}
	t.Execute(s, p)
	return strings.NewReader(s.String())
}

func (p *Percentage) HandleResp(resp *response.Response) {
	aggs := &percentAggs{}
	json.Unmarshal(resp.Aggregations, aggs)
	part := aggs.Part.DocCount
	whole := resp.Total()
	fmt.Println(part, whole)
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
