package main

import (
	"strconv"

	"lib/crypt"
	"lib/http"
	"lib/parser"
	"lib/proxy"
	"lib/tool"
)

func main() {

	serverHost := tool.ServerConfig["server"].(string)
	serverPort := tool.ServerConfig["server_port"].(float64)

	crypter := crypt.Bitcrypt{Secret: 0xB2}
	parser := parser.HTTPParser{}

	host := http.Host{
		Domain: serverHost,
		Port:   strconv.FormatFloat(serverPort, 'f', -1, 64),
	}

	server := proxy.Server{
		Host:    host,
		Crypter: crypter,
		Parser:  parser,
	}

	server.Start()
}
