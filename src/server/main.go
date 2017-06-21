package main

import (
	"strconv"
	"time"

	"lib/compressor"
	"lib/crypter"
	"lib/http"
	"lib/parser"
	"lib/proxy"
	"lib/tool"
)

func main() {

	serverHost := tool.ServerConfig["server"].(string)
	serverPort := tool.ServerConfig["server_port"].(float64)

	crypter := &crypter.Bitcrypter{Secret: 0xB2}
	parser := parser.HTTPParser{}
	compressor := &compressor.FlateCompressor{}

	host := http.Host{
		Domain: serverHost,
		Port:   strconv.FormatFloat(serverPort, 'f', -1, 64),
	}

	server := proxy.Server{
		Host:       host,
		Crypter:    crypter,
		Parser:     parser,
		Compressor: compressor,
		Deadline:   30 * time.Second,
	}

	server.Start()
}
