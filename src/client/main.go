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

	remote := tool.ClientConfig["server"].(string)
	remotePort := tool.ClientConfig["server_port"].(float64)
	local := tool.ClientConfig["local"].(string)
	localPort := tool.ClientConfig["local_port"].(float64)

	proxyHost := http.Host{Domain: remote, Port: strconv.FormatFloat(remotePort, 'f', -1, 64)}
	listen := http.Host{Domain: local, Port: strconv.FormatFloat(localPort, 'f', -1, 64)}
	client := proxy.Client{ProxyHost: proxyHost, Listen: listen, Crypter: crypter}
	client.Start()
}
