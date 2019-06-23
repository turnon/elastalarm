package response

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type formator interface {
	SetDetail([]interface{}, int, *big.Float)
	SetAbstract(string)
	String() string
}

func GetFormator(name string) func() formator {
	return func() formator {
		return &JSON{}
	}
}

type JSON struct {
	detailPairs []string
	abstract    string
}

func (j *JSON) SetDetail(arr []interface{}, count int, calculated *big.Float) {
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
