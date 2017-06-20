package main

import (
	"compress/flate"
	"fmt"
	"io"
	"log"
	"net"
)

//A a
type A interface {
	Foo()
}

//B b
type B struct {
}

// Foo foo 实现A接口
func (b B) Foo() {
	fmt.Println(b)
}

// Bar b 自有的方法
func (b B) Bar() {
	fmt.Println(b)
}

// Call call function
func Call(b A) {
	b.Foo()
	b.Bar() //报错,因为A接口没有Bar()方法
}

func main() {

	b := B{}
	Call(b)

	client, _ := net.Dial("tcp", "127.0.0.1:8900")
	c, err := flate.NewWriter(client, flate.DefaultCompression)
	if err != nil {
		log.Println("写入压缩流失败")
	}

	s := "This question continues the discussion started here. I found out that the HTTP response body can't be unmarshaled into JSON object because of deflate compression of the latter. Now I wonder how can I perform decompression with Golang. I will appreciate anyone who can show the errors in my code."
	io.WriteString(c, s)
	c.Flush()
}
