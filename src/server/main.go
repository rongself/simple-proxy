package main

import (
	"log"
	"strconv"
	"time"

	"lib/compressor"
	"lib/crypter"
	"lib/http"
	"lib/parser"
	"lib/proxy"
	"lib/tool"
	"strings"
)

func main() {

	serverHost := tool.ServerConfig["server"].(string)
	serverPort := tool.ServerConfig["server_port"].(float64)
	compress := tool.ServerConfig["compress"].(string)

	crypter := &crypter.Bitcrypter{Secret: 0xB2}
	parser := parser.HTTPParser{}

	host := http.Host{
		Domain: serverHost,
		Port:   strconv.FormatFloat(serverPort, 'f', -1, 64),
	}

	server := proxy.Server{
		Host:     host,
		Crypter:  crypter,
		Parser:   parser,
		Deadline: 2 * time.Hour,
	}

	if strings.Compare(compress, "") != 0 {
		compressor := compressor.FlateCompressor{}
		server.Compressor = &compressor
		log.Println("流量压缩开启")
	}

	server.Start()
}
