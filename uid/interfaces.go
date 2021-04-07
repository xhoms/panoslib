package uid

import "net/http"

type operation int

const (
	Login operation = iota
	Logout
	Group
	Ungroup
	Register
	Unregister
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type Monitor interface {
	Log(op operation, subject, value string, tout *uint)
}
