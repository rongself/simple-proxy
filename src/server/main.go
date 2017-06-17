package main

import (
	"strconv"

	"lib/crypt"
	"lib/http"
	"lib/proxy"
	"lib/tool"
)

func main() {

	// crypter := crypt.Bitcrypt{byte(rand.Intn(255))}
	crypter := crypt.Bitcrypt{Secret: 0xB2}
	serverHost := tool.ServerConfig["server"].(string)
	serverPort := tool.ServerConfig["server_port"].(float64)

	host := http.Host{Domain: serverHost, Port: strconv.FormatFloat(serverPort, 'f', -1, 64)}
	server := proxy.Server{Host: host, Crypter: crypter}
	server.Start()
}
