package main

import (
	"strconv"

	"lib/http"
	"lib/proxy"
	"lib/tool"
	"time"
)

func main() {

	remote := tool.ClientConfig["server"].(string)
	remotePort := tool.ClientConfig["server_port"].(float64)
	local := tool.ClientConfig["local"].(string)
	localPort := tool.ClientConfig["local_port"].(float64)
	compress := tool.ClientConfig["compress"].(string)
	// crypter := tool.ClientConfig["crypter"].(string)
	//password := tool.ClientConfig["password"].(string)
	// crypter := &crypter.Bitcrypter{Secret: 0xB2}

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
		Deadline:   2 * time.Hour,
		Compressor: compress,
	}

	client.Start()
}
