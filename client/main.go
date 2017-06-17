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

	remote := config.ClientConfig["server"].(string)
	remotePort := config.ClientConfig["server_port"].(float64)
	local := config.ClientConfig["local"].(string)
	localPort := config.ClientConfig["local_port"].(float64)

	proxyHost := http.Host{remote, strconv.FormatFloat(remotePort, 'f', -1, 64)}
	listen := http.Host{local, strconv.FormatFloat(localPort, 'f', -1, 64)}
	client := proxy.Client{proxyHost, listen, crypter}
	client.Start()
}
