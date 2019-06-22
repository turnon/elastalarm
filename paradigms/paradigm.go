package paradigms

import (
	"github.com/turnon/elastalarm/response"
)

type Paradigm interface {
	Template() string
	Found(resp *response.Response) (bool, *string)
	FoundOnAggs(resp *response.Response) (bool, *string)
	OnAggs() bool
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
