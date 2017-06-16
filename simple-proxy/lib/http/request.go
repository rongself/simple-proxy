package http

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
)

// Host host
type Host struct {
	Domain string
	Port   string
}

func (host Host) String() string {
	return fmt.Sprintf("%s:%s", host.Domain, host.Port)
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

// Request http request
type Request struct {
	Host     Host
	Method   string
	Protocol string
}

func (request Request) String() string {
	return fmt.Sprintf("%v %s %s", request.Host, request.Method, request.Protocol)
}

// Parse http string
func (request *Request) Parse(httpSteam []byte) error {

	var method, host, protocol string
	firstLine := string(httpSteam[:bytes.IndexByte(httpSteam[:], '\n')])

	reg := regexp.MustCompile(`^(GET|POST|DELETE|PUT|CONNECT|TRACE|PATCH|HEAD)\s.*$`)

	//丢弃不规范的请求
	if len(reg.FindStringIndex(firstLine)) <= 0 {
		return errors.New("请求不规范: " + firstLine)
	}

	fmt.Sscanf(firstLine, "%s%s%s", &method, &host, &protocol)

	var re = regexp.MustCompile(`^(\w+\:\/\/)`)
	if len(re.FindStringIndex(host)) <= 0 {
		if strings.IndexAny(host, ":443") > 0 {
			host = "https://" + host
		} else {
			host = "http://" + host
		}
	}

	u, err := url.Parse(host)
	if err != nil {
		log.Println("url解析失败", err)
	}

	domain := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "80"
	}
	request.Host.Domain = domain
	request.Host.Port = port
	request.Method = method
	request.Protocol = protocol
	// fmt.Println(request)
	return nil
}
