package tool

import (
	"log"
)

// Handle 处理错误
func Handle(err error, v ...interface{}) {
	if err != nil {
		log.Println(v, err)
	}
}

// HandleAndPanic 处理错误
func HandleAndPanic(err error, v ...interface{}) {
	if err != nil {
		log.Panicln(v, err)
	}
}

// Panic 处理错误
func Panic(err error) {
	if err != nil {
		panic(err)
	}
}
