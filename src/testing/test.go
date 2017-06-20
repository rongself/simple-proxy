package main

import (
	"compress/flate"
	"fmt"
	"io"
	"net"
	"time"
)

// func main() {

// 	host, _ := net.Listen("tcp", "127.0.0.1:8900")
// 	for {
// 		conn, err := host.Accept()
// 		if err != nil {
// 			log.Println("接收连接时出错", err)
// 		}
// 		fmt.Println("接收到一个连接")

// 		go handle(conn)

// 	}
// }
func main() {
	// 创建一个channel用以同步goroutine
	done := make(chan int, 2)

	// 在goroutine中执行输出操作
	go func() {

		// 告诉main函数执行完毕.
		// 这个channel在goroutine中是可见的
		// 因为它是在相同的地址空间执行的.
		for c := 0; c < 3; c++ {
			println("goroutine message")
			done <- c
		}

	}()

	time.Sleep(time.Second * 3)
	println("main function message")
	// 等待goroutine结束
	println(<-done, <-done, <-done)
	println("done!!")
}

func handle(conn net.Conn) {
	buffer := make([]byte, 1024)
	fmt.Println("等待发送数据")
	c := flate.NewReader(conn)
	for {
		n, err := c.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("break:", err)
				break
			}
			fmt.Println("错误", err)
			break
		}

		fmt.Println(string(buffer[:n]), buffer[:n], "--", n, len(buffer), cap(buffer))

	}
	// c.Close()}
}
