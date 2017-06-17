package main

import (
	"strconv"

	"../config"
	"../lib/crypt"
	"../lib/http"
	"../lib/proxy"
)

func main() {

	// crypter := crypt.Bitcrypt{byte(rand.Intn(255))}
	crypter := crypt.Bitcrypt{0xB2}
	serverHost := config.ServerConfig["server"].(string)
	serverPort := config.ServerConfig["server_port"].(float64)

	host := http.Host{serverHost, strconv.FormatFloat(serverPort, 'f', -1, 64)}
	server := proxy.Server{host, crypter}
	server.Start()
}
