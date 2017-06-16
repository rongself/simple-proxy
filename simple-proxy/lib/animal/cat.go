package cat

import (
	"fmt"
	"time"
)

//Cat cat
type Cat struct {
	Name      string
	Weight    float32
	IsRunning bool
}

func (cat Cat) run() {
	cat.IsRunning = true
	fmt.Printf("%s is running\n", cat.Name)
}

func (cat Cat) String() string {
	return fmt.Sprintf("{%s is %f weight}", cat.Name, cat.Weight)
}

func (cat *Cat) walk(seconds time.Duration) chan string {
	c := make(chan string)

	go func() {
		time.Sleep(seconds * time.Second)
		message := fmt.Sprintf("%s walk %d seconds\n", cat.Name, seconds)
		c <- message
	}()

	return c
}

// func main() {
// 	p := proxy.Server{}
// 	catTom := Cat{"Tom", 22.2, false}
// 	catKity := Cat{"Kity", 11.2, false}

// 	boom := catTom.walk(5)
// 	boom2 := catKity.walk(3)

// 	m1, m2 := <-boom, <-boom2

// 	fmt.Println(m1, m2, p)

// }
