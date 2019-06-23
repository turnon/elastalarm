package response

import "math/big"

type Result struct {
	abstract string
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

func (rs *Result) SetAbstract(abstract string) {
	rs.abstract = abstract
}

func (rs *Result) Stringify() string {
	f := GetFormator("")()
	for _, d := range rs.details {
		f.SetDetail(d.arr, d.raw, d.calculated)
	}
	f.SetAbstract(rs.abstract)
	return f.String()
}
