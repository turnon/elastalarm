package response

type Result struct {
	abstract string
	details  []*detail
}

type detail struct {
	arr        []interface{}
	raw        int
	calculated int
}

func (rs *Result) SetDetail(arr []interface{}, raw, calculated int) {
	d := &detail{arr, raw, calculated}
	rs.details = append(rs.details, d)
}

func (rs *Result) SetAbstract(abstract string) {
	rs.abstract = abstract
}

func (rs *Result) Stringify() string {
	f := GetFormator("")()
	for _, d := range rs.details {
		f.SetDetail(d.arr, d.raw)
	}
	f.SetAbstract(rs.abstract)
	return f.String()
}
