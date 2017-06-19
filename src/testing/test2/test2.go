package main

import (
	"compress/flate"
	"io"
	"log"
	"net"
)

func main() {
	client, _ := net.Dial("tcp", "127.0.0.1:8900")
	Compress
	c, err := flate.NewWriter(client, flate.DefaultCompression)
	if err != nil {
		log.Println("写入压缩流失败")
	}

	s := "This question continues the discussion started here. I found out that the HTTP response body can't be unmarshaled into JSON object because of deflate compression of the latter. Now I wonder how can I perform decompression with Golang. I will appreciate anyone who can show the errors in my code."
	io.WriteString(c, s)
	c.Flush()
}
