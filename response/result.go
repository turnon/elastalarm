package response

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type Result struct {
	Abstract string    `json:"abstract"`
	Details  []*detail `json:"details"`
}

type detail struct {
	Terms      []interface{} `json:"terms"`
	Raw        int           `json:"count"`
	Calculated *big.Float    `json:"calculated"`
}

func (rs *Result) SetDetail(arr []interface{}, raw int, calculated *big.Float) {
	d := &detail{arr, raw, calculated}
	rs.Details = append(rs.Details, d)
}

func (rs *Result) Text() string {
	return fmt.Sprintf("%s\n\n%s", rs.Abstract, rs.TextDetails())
}

func (rs *Result) TextDetails() string {
	var str strings.Builder

	for _, d := range rs.Details {
		str.WriteString(fmt.Sprint(d.termsStr()))
		str.WriteString(" ")
		RawStr := strconv.Itoa(d.Raw)
		if d.Calculated != nil {
			str.WriteString(d.Calculated.String())
			str.WriteString(" (")
			str.WriteString(RawStr)
			str.WriteString(")")
		} else {
			str.WriteString(RawStr)
		}
		str.WriteString("\n")
	}

	return str.String()
}

func (d *detail) termsStr() string {
	str, _ := json.Marshal(d.Terms)
	return string(str)
}
