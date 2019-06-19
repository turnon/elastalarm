package response

import (
	"fmt"
	"strconv"
	"strings"
)

type formator interface {
	SetDetail(arr []interface{}, count int)
	SetAbstract(abstract string)
	String() string
}

func GetFormator(name string) func() formator {
	return func() formator {
		if name == "json" {
			return &JSON{}
		}
		return &PlainText{}
	}
}

type PlainText struct {
	body     strings.Builder
	abstract string
}

func (t *PlainText) SetDetail(arr []interface{}, count int) {
	t.body.WriteString(fmt.Sprint(arr))
	t.body.WriteString(" ")
	t.body.WriteString(strconv.Itoa(count))
	t.body.WriteString("\n")
}

func (t *PlainText) SetAbstract(abstract string) {
	t.abstract = abstract
}

func (t *PlainText) String() string {
	return fmt.Sprintf("%s\n\n%s", t.abstract, t.body.String())
}

type JSON struct {
	detailPairs []string
	abstract    string
}

func (j *JSON) SetDetail(arr []interface{}, count int) {
	keys := []string{}
	for _, e := range arr {
		keys = append(keys, fmt.Sprint(e))
	}
	key := strings.Join(keys, " ")

	var detailPair strings.Builder
	detailPair.WriteString(`"`)
	detailPair.WriteString(key)
	detailPair.WriteString(`": `)
	detailPair.WriteString(strconv.Itoa(count))

	j.detailPairs = append(j.detailPairs, detailPair.String())
}

func (j *JSON) SetAbstract(abstract string) {
	j.abstract = abstract
}

func (j *JSON) String() string {
	var result strings.Builder
	result.WriteString(`{abstract: "`)
	result.WriteString(j.abstract)
	result.WriteString(`", detail:{`)

	for _, dp := range j.detailPairs {
		result.WriteString(dp)
	}

	result.WriteString(`}`)
	return result.String()
}
