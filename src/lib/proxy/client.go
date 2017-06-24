package proxy

import (
	"compress/flate"
	"io"
	"log"
	"net"
	"sync"

	"lib/compressor"
	"lib/conn"
	"lib/crypter"
	"lib/http"
	"lib/tool"
	"time"
)

const (
	//ClientReadTimeOut 客户端连接超时
	ClientReadTimeOut = 5 * time.Hour
	//ClientWriteTimeOut 客户端写超时
	ClientWriteTimeOut = 5 * time.Hour
	//MaxProxyConn 允许最大Porxy服务器连接数
	MaxProxyConn = 2
)

//Client proxy client
type Client struct {
	ProxyHost  http.Host
	Listen     http.Host
	Crypter    crypter.Crypter
	Compressor compressor.Compressor
	Deadline   time.Duration
}

// ConnProxy 连接Proxy服务器
func (client Client) ConnProxy(ip *net.TCPAddr) (io.ReadWriteCloser, error) {

	proxyServer, err := net.DialTCP("tcp", nil, ip)
	tool.HandleAndPanic(err, "连接Proxy失败")

	//压缩选项
	if client.Compressor != nil {
		err := client.Compressor.Init(proxyServer, flate.DefaultCompression)
		tool.HandleAndPanic(err, "初始化压缩器失")
		return client.Compressor, err
	}

	return proxyServer, err

}

// Start start proxy client
func (client Client) Start() {

	localServer, err := net.Listen("tcp", client.Listen.String())
	tool.HandleAndPanic(err, "端口监听失败")

	log.Println("客户端开始监听端口:", client.Listen.String())
	defer localServer.Close()

	pool, err := conn.InitPool(MaxProxyConn, 2)
	tool.HandleAndPanic(err, "初始化连接池错误")

	ip, err := net.ResolveTCPAddr("tcp", client.ProxyHost.String())
	tool.HandleAndPanic(err, "IP解析失败: ")

	// 为连接池填充连接
	for i := 0; i < MaxProxyConn; i++ {

		proxyServer, err := client.ConnProxy(ip)
		tool.HandleAndPanic(err, "连接Proxy服务器失败: ")
		log.Println("worker", proxyServer, "连接Proxy服务器成功:", client.ProxyHost.String())
		pool.Push(proxyServer)
		defer func() {
			proxyServer.Close()
			log.Println("worker", proxyServer, "关闭")
		}()

	}

	for {
		brower, err := localServer.Accept()
		if err != nil {
			log.Panic("接受连接失败", err)
		}
		log.Println("浏览器连接成功:", brower.RemoteAddr().String())
		brower.SetDeadline(time.Now().Add(client.Deadline))
		now := time.Now()

		go func() {
			client.HandleRequest(brower, &pool)
			log.Println("请求处理完成,处理时间:", time.Since(now))
		}()
	}
}

// HandleRequest handle
func (client Client) HandleRequest(brower net.Conn, pool *conn.Pool) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	defer brower.Close()

	// 写在handle里,每一个浏览器请求建立一个TCP连接来传送
	proxyServer := pool.Pop().(io.ReadWriteCloser)
	tool.Debug("取出,剩余:", pool.Len())
	//代理过来的流量写回到浏览器

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		w, err := tool.Copy(brower, proxyServer, client.Crypter)
		log.Println("client -> brower", w)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				log.Println("超时:client -> brower", err)
			} else {
				log.Println("错误:client -> brower", w, err)
			}
			defer proxyServer.Close()
			ip, err := net.ResolveTCPAddr("tcp", client.ProxyHost.String())
			tool.HandleAndPanic(err, "IP解析失败: ")
			refeed, err := client.ConnProxy(ip)
			tool.HandleAndPanic(err, "补充Proxy服务器失败: ")
			pool.Push(refeed)
			tool.Debug("异常归还", pool.Len())
		}
		wg.Done()
	}()

	//浏览器过来的流量写入到代理服务器
	wg.Add(1)
	go func() {
		w, err := tool.Copy(proxyServer, brower, client.Crypter)
		log.Println("brower -> proxy", w)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				log.Println("超时:brower -> proxy", err)
			} else {
				log.Println("错误:brower -> proxy", w, err)
			}
			defer proxyServer.Close()

			ip, err := net.ResolveTCPAddr("tcp", client.ProxyHost.String())
			tool.HandleAndPanic(err, "IP解析失败: ")
			refeed, err := client.ConnProxy(ip)
			tool.HandleAndPanic(err, "补充Proxy服务器失败: ")
			pool.Push(refeed)
			tool.Debug("异常归还", pool.Len())

		} else {
			pool.Push(proxyServer)
			tool.Debug("正常归还2", pool.Len())
		}
		wg.Done()
	}()

	wg.Wait()
	log.Println("done")

}
