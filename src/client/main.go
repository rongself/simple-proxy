package main

import (
	"strconv"

	"lib/compressor"
	"lib/crypter"
	"lib/http"
	"lib/proxy"
	"lib/tool"
)

func main() {

	remote := tool.ClientConfig["server"].(string)
	remotePort := tool.ClientConfig["server_port"].(float64)
	local := tool.ClientConfig["local"].(string)
	localPort := tool.ClientConfig["local_port"].(float64)

	crypter := crypter.Bitcrypter{Secret: 0xB2}
	compressor := compressor.FlateCompressor{}

	proxyHost := http.Host{
		Domain: remote,
		Port:   strconv.FormatFloat(remotePort, 'f', -1, 64),
	}

	listen := http.Host{
		Domain: local,
		Port:   strconv.FormatFloat(localPort, 'f', -1, 64),
	}

	client := proxy.Client{
		ProxyHost:  proxyHost,
		Listen:     listen,
		Crypter:    crypter,
		Compressor: &compressor,
	}

	client.Start()
}
