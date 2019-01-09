package paradigms

import (
	"math/big"

	"github.com/turnon/elastalarm/response"
)

type Paradigm interface {
	Template() string
	Found(resp *response.Response) (bool, *string)
}

func Names(name string) Paradigm {
	switch name {
	case "percentage":
		return &Percentage{}
	case "count":
		return &Count{}
	case "spike":
		return &Spike{}
	default:
		return nil
	}
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
