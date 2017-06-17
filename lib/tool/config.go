package tool

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//JSONReader reader
type JSONReader struct {
	file   []byte
	schema interface{}
}

// NewJSONReader 构造函数
func NewJSONReader(file string, schema interface{}) JSONReader {

	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Panic("配置文件不存在: ", file, err)
	}
	return JSONReader{fileBytes, schema}
}

func (config JSONReader) Read() interface{} {
	json.Unmarshal(config.file, &(config.schema))
	return config.schema
}
