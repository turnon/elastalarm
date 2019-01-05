package paradigms

import (
	"math/big"

	"bitbucket.org/xcrossing/elastic_alarm/response"
)

type Paradigm interface {
	Template() string
	Found(resp *response.Response) (bool, *string)
}

type Match struct {
	Gt, Lt *float64
}

func (m *Match) ing(v *big.Float) bool {
	if m.Gt != nil && v.Cmp(configValue(m.Gt)) != 1 {
		return false
	}

	if m.Lt != nil && v.Cmp(configValue(m.Lt)) != -1 {
		return false
	}

	return true
}

func configValue(v *float64) *big.Float {
	return big.NewFloat(*v)
}
