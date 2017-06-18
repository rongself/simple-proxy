package parser

import "lib/http"

// Parser 解析器接口
type Parser interface {
	Parse(httpSteam []byte) (http.Request, error)
}
