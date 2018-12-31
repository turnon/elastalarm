package paradigms

import "io"

type Paradigm interface {
	ReqBody() io.Reader
}
