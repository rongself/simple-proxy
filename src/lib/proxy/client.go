package proxy

import (
	"compress/flate"
	"io"
	"log"
	"net"
	"sync"

	"lib/compressor"
	"lib/http"
	"lib/tool"
	"time"
)

//Client proxy client
type Client struct {
	ProxyHost   http.Host
	Listen      http.Host
	Compressor  string
	Deadline    time.Duration
	CheckOnline time.Duration
}

// ConnProxy 连接Proxy服务器
func (client Client) ConnProxy(ip *net.TCPAddr) (io.ReadWriteCloser, error) {

	proxyServer, err := net.DialTCP("tcp", nil, ip)
	tool.HandleAndPanic(err, "连接Proxy失败")
	proxyServer.SetKeepAlive(true)
	proxyServer.SetKeepAlivePeriod(client.CheckOnline)
	//压缩选项
	if client.Compressor != "" {
		comp, err := compressor.NewCompressor(proxyServer, flate.DefaultCompression)
		tool.HandleAndPanic(err, "初始化压缩器失败")
		return comp, err
	}
	return proxyServer, err
}

// Start start proxy client
func (client Client) Start() {

	ip, err := net.ResolveTCPAddr("tcp", client.Listen.String())
	tool.HandleAndPanic(err, "IP解析失败: ")
	localServer, err := net.ListenTCP("tcp", ip)
	tool.HandleAndPanic(err, "端口监听失败")

	log.Println("客户端开始监听端口:", client.Listen.String())
	defer localServer.Close()

	for {
		brower, err := localServer.Accept()
		if err != nil {
			log.Panic("接受连接失败", err)
		}
		log.Println("浏览器连接成功:", brower.RemoteAddr().String())
		if browerTCP, ok := brower.(*net.TCPConn); ok {
			now := time.Now()
			go func() {
				err := browerTCP.SetKeepAlive(true)
				tool.Handle(err, "SetKeepAlive failed")
				err = browerTCP.SetKeepAlivePeriod(client.CheckOnline)
				tool.Handle(err, "SetKeepAlivePeriod failed")
				client.HandleRequest(browerTCP)
				log.Println("请求处理完成,处理时间:", time.Since(now))
			}()
		}
		//brower.SetDeadline(time.Now().Add(client.Deadline))
	}
}

// HandleRequest handle
func (client Client) HandleRequest(brower *net.TCPConn) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	defer brower.Close()

	// 写在handle里,每一个浏览器请求建立一个TCP连接来传送
	ip, err := net.ResolveTCPAddr("tcp", client.ProxyHost.String())
	tool.Handle(err, "IP解析失败: ")
	proxyServer, err := client.ConnProxy(ip)
	tool.Handle(err, "连接Proxy服务器失败: ")
	defer proxyServer.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	//代理过来的流量写回到浏览器
	go func() {
		w, err := tool.Copy(brower, proxyServer)
		log.Println("client -> brower", w)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				log.Println("超时:client -> brower", err)
			} else {
				log.Println("错误:client -> brower", w, err)
			}
		}
		wg.Done()
	}()

	//浏览器过来的流量写入到代理服务器
	wg.Add(1)
	go func() {
		w, err := tool.Copy(proxyServer, brower)
		log.Println("brower -> proxy", w)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				log.Println("超时:brower -> proxy", err)
			} else {
				log.Println("错误:brower -> proxy", w, err)
			}
		}
		wg.Done()
	}()

	wg.Wait()
	log.Println("done")

}
