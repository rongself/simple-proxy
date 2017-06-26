package conn

import (
	"errors"
	"net"
)

// var pool = map[string][]int{

// 	"127.0.0.1:8765": []int{1, 2, 3},
// }

//Conn conn
type Conn net.TCPConn

func (conn *Conn) isClosed() bool {

	buffer := make([]byte, 1)
	_, err := conn.Read(buffer)
	return err != nil
}

//Pool pool
type Pool struct {
	Conns       chan net.Conn
	MaxConn     int
	MaxIdleConn int
}

//InitPool init a pool
func InitPool(maxConn int, maxIdleConn int) (Pool, error) {
	var err error
	if maxConn < maxIdleConn {
		err = errors.New("maxIdleConn不能大于maxIdleConn")
	}
	pool := Pool{
		Conns:       make(chan net.Conn, maxConn),
		MaxConn:     maxConn,
		MaxIdleConn: maxIdleConn,
	}
	return pool, err
}

//Len length
func (pool Pool) Len() int {
	return len(pool.Conns)
}

//Pop pop
func (pool *Pool) Pop() interface{} {
	return <-pool.Conns
}

//Push push
func (pool *Pool) Push(conn interface{}) {
	(*pool).Conns <- conn.(*net.TCPConn)
}

// Pools pool
type Pools struct {
	pools map[string][]Pool
}

func xx() {
	// p := Pool{}
	// heap.Pop(&p)
}
