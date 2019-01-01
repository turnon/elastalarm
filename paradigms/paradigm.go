package paradigms

import (
	"io"

	"bitbucket.org/xcrossing/elastic_alarm/response"
)

type Paradigm interface {
	ReqBody() io.Reader
	HandleResp(resp *response.Response)
}
