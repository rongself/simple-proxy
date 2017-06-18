package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	// b := make([]byte, 5)
	host, _ := net.Listen("tcp", "127.0.0.1:8900")
	for {
		conn, err := host.Accept()
		if err != nil {
			log.Println("接收连接时出错", err)
		}
		fmt.Println("接收到一个连接")

		go handle(conn)

	}
}

func handle(conn net.Conn) {
	for {

		buffer := make([]byte, 5)
		fmt.Println("等待发送数据")
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		}
		if bytes.Contains(buffer, []byte{'\r', '\n'}) {
			fmt.Println("\\n find")
		}
		fmt.Println(buffer, "--", n, len(buffer), cap(buffer))
	}
}
