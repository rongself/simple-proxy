package tool

import (
	"log"
)

// DebugOutput Debug信息输出开关
var DebugOutput = true

//Debug 输出debug日志
func Debug(v ...interface{}) {
	if DebugOutput {
		log.Println(v)
	}
}

//Info 输出日志
func Info(v ...interface{}) {
	log.Println(v)
}

//Error 输出错误信息
func Error(v ...interface{}) {

}
