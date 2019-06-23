package response

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type Result struct {
	Abstract string
	details  []*detail
}

type detail struct {
	arr        []interface{}
	raw        int
	calculated *big.Float
}

func (rs *Result) SetDetail(arr []interface{}, raw int, calculated *big.Float) {
	d := &detail{arr, raw, calculated}
	rs.details = append(rs.details, d)
}

func (rs *Result) Text() string {
	return fmt.Sprintf("%s\n\n%s", rs.Abstract, rs.TextDetails())
}

func (rs *Result) TextDetails() string {
	var str strings.Builder

	for _, d := range rs.details {
		str.WriteString(fmt.Sprint(d.arr))
		str.WriteString(" ")
		if d.calculated != nil {
			str.WriteString(d.calculated.String())
			str.WriteString(" (")
			str.WriteString(strconv.Itoa(d.raw))
			str.WriteString(")")
		} else {
			str.WriteString(strconv.Itoa(d.raw))
		}
		str.WriteString("\n")
	}

	return str.String()
}
