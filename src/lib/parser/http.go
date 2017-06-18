package parser

import (
	"bytes"
	"errors"
	"fmt"
	"lib/http"
	"log"
	"net/url"
	"regexp"
	"strings"
)

// HTTPParser http解析器
type HTTPParser struct {
}

// Parse http string
func (parser HTTPParser) Parse(httpSteam []byte) (http.Request, error) {

	var request = http.Request{}
	var method, host, protocol string
	firstLine := string(httpSteam[:bytes.IndexByte(httpSteam[:], '\n')])

	reg := regexp.MustCompile(`^(GET|POST|DELETE|PUT|CONNECT|TRACE|PATCH|HEAD)\s.*$`)

	//丢弃不规范的请求
	if len(reg.FindStringIndex(firstLine)) <= 0 {
		return request, errors.New("请求不规范: " + firstLine)
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
	return request, nil
}
