package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

var conns []*net.Conn

func main() {
	host, _ := net.Listen("tcp", "127.0.0.1:8900")
	for {
		conn, err := host.Accept()
		if err != nil {
			log.Println("接收连接时出错", err)
		}
		fmt.Println("接收到一个连接来自", conn.RemoteAddr())
		conns = append(conns, &conn)
		go handle(&conn)

	}
}

func handle(conn *net.Conn) {
	buffer := make([]byte, 512)
	fmt.Println("等待发送数据", conns)
	broadcastString(conn, fmt.Sprintf("%v上线了\n", (*conn).RemoteAddr()))
	// c := flate.NewReader(conn)
	for {
		n, err := (*conn).Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("IO EOF:", err, buffer[:n])
				broadcastString(conn, fmt.Sprintf("%v下线了\n", (*conn).RemoteAddr()))
				break
			}
			fmt.Println("错误", err)
			break
		}

		broadcast(conn, buffer[:n])

	}
	defer func() {
		(*conn).Close()
		fmt.Println("连接接关闭", (*conn).RemoteAddr())
	}()
}

func broadcast(sender *net.Conn, messge []byte) {
	var wg sync.WaitGroup
	for i := 0; i < len(conns); i++ {
		if conns[i] != sender {
			wg.Add(1)
			go func(index int) {
				connt := *(conns[index])
				connt.Write(messge)
				wg.Done()
			}(i)

		}
	}
	wg.Wait()
}

func broadcastString(sender *net.Conn, messge string) {
	broadcast(sender, []byte(messge))
}
