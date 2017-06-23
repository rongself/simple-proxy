package conn

import (
	"container/heap"
	"errors"
	"net"
	"sync"
)

// var pool = map[string][]int{

// 	"127.0.0.1:8765": []int{1, 2, 3},
// }

const (
	// MaxConn max conn
	MaxConn = 9
)

//Conn conn
type Conn net.TCPConn

func (conn *Conn) isClosed() bool {

	buffer := make([]byte, 1)
	_, err := conn.Read(buffer)
	return err != nil
}

//Pool pool
type Pool struct {
	Conns       []net.Conn
	MaxConn     int
	MaxIdleConn int
	lock        *sync.Mutex
}

//InitPool init a pool
func InitPool(maxConn int, maxIdleConn int) (Pool, error) {
	var err error
	if maxConn < maxIdleConn {
		err = errors.New("maxIdleConn不能大于maxIdleConn")
	}
	pool := Pool{
		MaxConn:     maxConn,
		MaxIdleConn: maxIdleConn,
		lock:        new(sync.Mutex),
	}
	return pool, err
}

//Len length
func (pool Pool) Len() int {
	return len(pool.Conns)
}

//Less less
func (pool Pool) Less(i, j int) bool {
	return true
}

//Swap swap
func (pool Pool) Swap(i, j int) {
	pool.Conns[i], pool.Conns[j] = pool.Conns[j], pool.Conns[i]
}

//Pop pop
func (pool *Pool) Pop() interface{} {
	conns := (*pool).Conns
	n := len(conns)
	conn := conns[n-1]
	(*pool).Conns = conns[0 : n-1]
	return conn
}

//Push push
func (pool *Pool) Push(conn interface{}) {
	(*pool).Conns = append((*pool).Conns, conn.(*net.TCPConn))
}

//SyncPop 同步pop操作
func (pool *Pool) SyncPop() interface{} {
	var item interface{}
	for {
		if pool.Len() <= 0 {
			continue
		}
		if item = heap.Pop(pool); item != nil {
			return item
		}
	}
}

//LockPop 带锁的pop
func (pool *Pool) LockPop() interface{} {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	return heap.Pop(pool)
}

// LockPush 带锁的Push
func (pool *Pool) LockPush(conn interface{}) {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	heap.Push(pool, conn)
}

// Pools pool
type Pools struct {
	pools map[string][]Pool
}

func xx() {
	p := Pool{}
	heap.Pop(&p)
}
