package paradigms

import (
	"github.com/turnon/elastalarm/response"
)

type Paradigm interface {
	Template() string
	Found(resp *response.Response) (bool, *response.Result)
	FoundOnAggs(resp *response.Response) (bool, *response.Result)
	OnAggs() bool
	SupportStep() bool
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
