package tool

// ServerConfig 服务器端配置格式
var ServerConfig = map[string]interface{}{
	"server":      "127.0.0.1",
	"server_port": 8388,
	"local_port":  1080,
	"password":    "barfoo!",
	"method":      "bitcrypt",
	"timeout":     60,
}

// ClientConfig 客户端端配置格式
var ClientConfig = map[string]interface{}{
	"server":      "127.0.0.1",
	"server_port": 8388,
	"local":       "0.0.0.0",
	"local_port":  1080,
	"password":    "barfoo!",
	"method":      "bitcrypt",
	"timeout":     60,
}

func init() {
	NewJSONReader("./config/server.config.json", &ServerConfig).Read()
	NewJSONReader("./config/client.config.json", &ClientConfig).Read()
}
