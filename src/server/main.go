package main

import (
	"strconv"
	"time"

	"lib/http"
	"lib/parser"
	"lib/proxy"
	"lib/tool"
)

func main() {

	serverHost := tool.ServerConfig["server"].(string)
	serverPort := tool.ServerConfig["server_port"].(float64)
	compress := tool.ServerConfig["compress"].(string)

	// crypter := &crypter.Bitcrypter{Secret: 0xB2}
	parser := parser.HTTPParser{}

	host := http.Host{
		Domain: serverHost,
		Port:   strconv.FormatFloat(serverPort, 'f', -1, 64),
	}

	server := proxy.Server{
		Host:       host,
		Parser:     parser,
		Deadline:   2 * time.Hour,
		Compressor: compress,
	}

	server.Start()
}
