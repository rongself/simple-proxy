# simple-proxy

## 简介
一个快速、轻量级的Proxy

使用Golang编写的一个简单的Proxy,采用CS结构,加密使用轻便的位运算加密,速度还不错,这是我写的第一个Golang程序,如有错漏迎提Iusse指出,也欢迎提交PR

## 环境

- Golang 1.8.x

## 安装

```bash
# clone 项目
git clone https://github.com/rongself/simple-proxy.git

# 编译项目
cd simple-proxy && ./install.sh
```
执行成功之后会在./bin文件夹下生成以下

```
./
├── bin
│   ├── client                 // 客户端可执行文件
│   ├── config
│   │   ├── client.config.json // 客户端配置文件
│   │   └── server.config.json // 服务器端配置文件
│   └── server                 // 服务器端可执行文件

```

## 配置

### 服务器端

编辑配置文件 `./bin/server.config.json` 

```

{
    "server":"0.0.0.0",     //服务器监听IP,一般设为 0.0.0.0
    "server_port":8888,     //服务器监听端口
    "password":"barfoo!",   //* 密码,此项暂时还没有实现
    "method": "bitcrypt",   //* 加密方式,此项暂时还没有实现,默认使用位运加密
    "timeout":60            //* 连接超时时间,此项暂时还没有实现
}

```

### 客户端

编辑配置文件 `./bin/client.config.json` 

```

{
    "server":"yourserver.com",      //服务器监听IP,一般设为 0.0.0.0
    "server_port":8888,             //服务器监听端口
    "local":"0.0.0.0",              //本地监听IP,设为127.0.0.1只允许本地连接,为0.0.0.0允许局域网链接
    "local_port":1070,              //本地监听端口,浏览器代理设置此端口
    "password":"barfoo!",           //* 密码,此项暂时还没有实现
    "method": "bitcrypt",           //* 加密方式,此项暂时还没有实现,默认使用位运加密
    "timeout":60                    //* 连接超时时间,此项暂时还没有实现
}
```

### 运行

运行二进制文件,`必须`在`./bin`文件夹中运行,否则无法读取到配置文件

1. 在外网服务器运行服务器端:

```bash
cd ./bin && ./server
```
2. 在本地运行客户端

```bash
cd ./bin && ./client
```

## 使用

设浏览器代理为客户端配置中配置的地址`[本地IP|局域网IP]:local_port`,上面配置即为 `127.0.0.1:1070`
