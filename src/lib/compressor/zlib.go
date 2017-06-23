package compressor

// import "io"
// import "compress/zlib"
// import "log"

// // ZlibCompressor Flate压缩器
// type ZlibCompressor struct {
// }

// // NewWriter 创建新的写入器
// func (compressor *ZlibCompressor) NewWriter(w io.Writer, level int) (WriteCloser, error) {
// 	return zlib.NewWriterLevel(w, level)
// }

// // NewReader 创建新的读取器
// func (compressor *ZlibCompressor) NewReader(r io.Reader) io.ReadCloser {
// 	zr, err := zlib.NewReader(r)
// 	if err != nil {
// 		log.Println("zlib error:", err)
// 	}
// 	return zr
// }
