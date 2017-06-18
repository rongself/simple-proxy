package http

import "fmt"

// Request http request
type Request struct {
	Host     Host
	Method   string
	Protocol string
}

// Host host
type Host struct {
	Domain string
	Port   string
}

// Http Methods
const (
	GET     = "GET"
	POST    = "POST"
	DELETE  = "DELETE"
	PUT     = "PUT"
	CONNECT = "CONNECT"
	TRACE   = "TRACE"
	PATCH   = "PATCH"
	HEAD    = "HEAD"
)

func (host Host) String() string {
	return fmt.Sprintf("%s:%s", host.Domain, host.Port)
}

func (request Request) String() string {
	return fmt.Sprintf("%v %s %s", request.Host, request.Method, request.Protocol)
}
