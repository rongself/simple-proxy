package proxy

import (
	"compress/flate"
	"io"
	"log"
	"net"

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
	MaxProxyConn = 5
)

//Client proxy client
type Client struct {
	ProxyHost  http.Host
	Listen     http.Host
	Crypter    crypter.Crypter
	Compressor compressor.Compressor
	Deadline   time.Duration
}

// Start start proxy client
func (client Client) Start() {

	localServer, err := net.Listen("tcp", client.Listen.String())
	if err != nil {
		log.Panic("端口监听失败", err)
	}

	log.Println("客户端开始监听端口:", client.Listen.String())
	defer localServer.Close()

	pool, err := conn.InitPool(MaxProxyConn, 5)
	tool.Handle(err, "初始化连接池错误")

	// ip, err := net.ResolveTCPAddr("tcp", client.ProxyHost.String())
	// if err != nil {
	// 	log.Panic("IP解析失败: ", err)
	// }

	// for i := 0; i < MaxProxyConn; i++ {

	// 	proxyServer, err := net.DialTCP("tcp", nil, ip)
	// 	if err != nil {
	// 		log.Panic("连接Proxy服务器失败: ", err)
	// 	}
	// 	log.Println("worker", proxyServer, "连接Proxy服务器成功:", client.ProxyHost.String())
	// 	defer func() {
	// 		proxyServer.Close()
	// 		log.Println("worker", proxyServer, "关闭")
	// 	}()

	// 	var cr io.ReadCloser
	// 	var cw io.WriteCloser
	// 	if client.Compressor != nil {
	// 		// flate.NewWriter()
	// 		var err error
	// 		cr = client.Compressor.NewReader(proxyServer)
	// 		cw, err = client.Compressor.NewWriter(proxyServer, flate.DefaultCompression)
	// 		if err != nil {
	// 			log.Panic("初始化压缩器失败", err)
	// 		}
	// 	} else {
	// 		cr, cw = proxyServer, proxyServer
	// 	}
	// 	defer cr.Close()
	// 	defer cw.Close()

	// 	pool.LockPush(proxyServer)

	// }

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

	ip, err := net.ResolveTCPAddr("tcp", client.ProxyHost.String())
	if err != nil {
		log.Panic("IP解析失败: ", err)
	}
	// 写在handle里,每一个浏览器请求建立一个TCP连接来传送
	proxyServer, err := net.DialTCP("tcp", nil, ip)
	if err != nil {
		log.Panic("连接Proxy服务器失败: ", err)
	}
	proxyServer.SetDeadline(time.Now().Add(client.Deadline))
	log.Println("连接Proxy服务器成功:", client.ProxyHost.String())
	defer proxyServer.Close()

	// 包装代理服务器通道

	var compressProxyConn io.ReadWriteCloser
	if client.Compressor != nil {
		// flate.NewWriter()
		var err error
		client.Compressor.Init(proxyServer, flate.DefaultCompression)
		compressProxyConn = client.Compressor
		if err != nil {
			log.Panic("初始化压缩器失败", err)
		}
	} else {
		compressProxyConn = proxyServer
	}
	defer compressProxyConn.Close()

	//代理过来的流量写回到浏览器
	channel := make(chan bool, 2)
	defer close(channel)
	go func() {
		w, err := tool.Copy(brower, compressProxyConn, client.Crypter)
		log.Println("client -> brower", w)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				log.Println("超时:client -> brower", err)
			} else {
				log.Println("错误:client -> brower", w, err)
			}
		}
		channel <- true
	}()

	//浏览器过来的流量写入到代理服务器
	go func() {
		w, err := tool.Copy(compressProxyConn, brower, client.Crypter)
		log.Println("brower -> proxy", w)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				log.Println("超时:brower -> proxy", err)
			} else {
				log.Println("错误:brower -> proxy", w, err)
			}
		}
		channel <- true
	}()

	w1, w2 := <-channel, <-channel
	log.Println("w1:", w1, "w2:", w2)

}
