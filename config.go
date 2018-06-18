package main

import (
	"github.com/biluohc/kloud/api"
	"github.com/go-pg/pg"
)

func init() {
	api.Host = Host
}

// 数据库的配置
var dbConfig = pg.Options{
	// TCP host:port or Unix socket depending on Network.
	Network:  "tcp",
	Addr:     "localhost:5432",
	User:     "pgmxo",
	Database: "pgmxodb",
	Password: "imissyou.",
}

// 服务器监听的地址
const ServeAddr string = "0.0.0.0:8080"

// 服务器在网络上的地址

// const Host string = "http://biluohc.me:8080"
const Host string = "http://localhost:8080/"

// 日志的配置在log.go和log开头的那几个xml

// 阿里云不备案无法使用 let's crypto.., 223, 手动再见
// https 相关
var https bool = false
var httpsCrt string = "biluohc.me.crt"
var httpsKey string = "biluohc.me.key"

// env GO_LOG=2 nohup  ./kloud  &
