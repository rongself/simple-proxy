package main

import (
	"compress/flate"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {

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
	buffer := make([]byte, 1024)
	fmt.Println("等待发送数据")
	c := flate.NewReader(conn)
	n, err := conn.Read(buffer)
	if err != nil {
		if err == io.EOF {
			fmt.Println("break:", err)
		}
		fmt.Println("错误", err)
	}
	// c.Close()
	fmt.Println(string(buffer[:n]), buffer[:n], "--", n, len(buffer), cap(buffer))
}
