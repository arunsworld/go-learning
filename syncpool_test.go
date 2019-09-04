package learning

import (
	"log"
	"sync"
	"testing"
)

func TestSyncPool(t *testing.T) {
	data := testSyncPool()
	for i := 0; i < 9999991; i++ {
		if data[i] != i {
			log.Fatal("Did not get expected number...")
		}
	}
}

func BenchmarkSyncPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testSyncPool()
	}
}

func testSyncPool() []int {
	pool := sync.Pool{
		New: func() interface{} {
			return make([]int, 1000)
		},
	}
	ch := make(chan []int, 1000)
	go func() {
		// buffer := make([]int, 1000)
		buffer := pool.Get().([]int)
		counter := 0
		for i := 0; i < 9999991; i++ {
			if counter == 1000 {
				ch <- buffer
				// buffer = make([]int, 1000)
				buffer = pool.Get().([]int)
				counter = 0
			}
			buffer[counter] = i
			counter++
		}
		ch <- buffer[:counter]
		close(ch)
	}()
	data := make([]int, 0, 300000000)
	// data := []int{}
	for buffer := range ch {
		// for _, v := range buffer {
		// 	data = append(data, v)
		// }
		data = append(data, buffer...)
		pool.Put(buffer)
	}
	return data
}
