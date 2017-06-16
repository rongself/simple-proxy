package main

import (
	"./lib/crypt"
	"./lib/http"
	"./lib/proxy"
)

func main() {
	crypter := crypt.Bitcrypt{}
	host := http.Host{"127.0.0.1", "8888"}
	server := proxy.Server{host, crypter}
	go server.Start()

	proxyHost := http.Host{"127.0.0.1", "8888"}
	listen := http.Host{"127.0.0.1", "1090"}
	client := proxy.Client{proxyHost, listen, crypter}
	client.Start()

	// var Secret = []byte{0xB2, 0x09, 0xBB, 0x55, 0x93, 0x6D, 0x44, 0x47}
	// fmt.Println(Secret)
	// fmt.Println(c.Encode(Secret))
	// fmt.Println(c.Decode(Secret))
}
