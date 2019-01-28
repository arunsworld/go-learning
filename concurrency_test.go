package learning

import (
	"sync"
	"testing"
)

func TestBasicConcurrency(t *testing.T) {
	ch := make(chan int)
	// Publish 10 messages to the channel in sequence
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()
	// Now read them and ensure we got in sequence
	for i := 0; i < 10; i++ {
		v := <-ch
		if i != v {
			t.Errorf("Expected %d got %d", i, v)
		}
	}
}

func TestRunningConcurrentJobsAndWaitingForThemToComplete(t *testing.T) {
	wg := sync.WaitGroup{}
	// Kick off 10 jobs
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			// Code for dummy job - after completion mark as Done (defer it in the event of early returns)
			defer wg.Done()
		}()
	}
	wg.Wait()
}

type Cache interface {
	Store(key int, value int)
	Retrieve(key int) int
}

func TestThreadSafeMapCache(t *testing.T) {
	var cache Cache
	cache = createCache()
	// Now store 10 values in parallel
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			cache.Store(v, v)
			wg.Done()
		}(i)
	}
	wg.Wait()
	// Now retreive the values in parallel and ensure we're good
	wg = sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			newV := cache.Retrieve(v)
			if v != newV {
				t.Errorf("Expected to get: %d. Got %d instead.", v, newV)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

/*
Below implementation of cache does not work because map is not thread-safe.
	type cache struct {
		store map[int]int
	}

	func (c *cache) Store(key int, value int) {
		c.store[key] = value
	}

	func (c *cache) Retrieve(key int) int {
		return c.store[key]
	}
*/

type cache struct {
	store sync.Map
}

func (c *cache) Store(key int, value int) {
	c.store.Store(key, value)
}

func (c *cache) Retrieve(key int) int {
	v, ok := c.store.Load(key)
	if !ok {
		return 0
	}
	return v.(int)
}

func createCache() Cache {
	result := &cache{}
	result.store = sync.Map{}
	return result
}

func BenchmarkThreadSafeMapCache(b *testing.B) {
	cache := createCache()
	for n := 0; n < b.N; n++ {
		cache.Store(5, 5)
	}
	for n := 0; n < b.N; n++ {
		cache.Retrieve(5)
	}
}

func TestPublishToMultipleListeners(t *testing.T) {
	listenerCh1 := make(chan int)
	listenerCh2 := make(chan int)
	ch := make(chan struct{})
	go func() {
		for i := 0; i < 1000; i++ {
			ch <- struct{}{}
		}
		close(ch)
	}()
	// Listen to ch and publish to listenerCh1
	go func() {
		counter := 0
		for {
			_, ok := <-ch
			if !ok {
				break
			}
			counter++
		}
		listenerCh1 <- counter
	}()
	// Listen to ch and publish to listenerCh2
	go func() {
		counter := 0
		for {
			_, ok := <-ch
			if !ok {
				break
			}
			counter++
		}
		listenerCh2 <- counter
	}()
	a := <-listenerCh1
	b := <-listenerCh2
	if a+b != 1000 {
		t.Fatal("Expected responses to total 1000 but didn't")
	}
}

func TestBroadcastUsingChannelClose(t *testing.T) {
	broadcast := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		// Kick off a job waiting on broadcast after which will mark wg as Done
		go func() {
			// Do some serious work
			<-broadcast
			wg.Done()
		}()
	}
	// Now broadcast to jobs by closing the channel
	go func() {
		// Do some serious work
		close(broadcast)
	}()
	wg.Wait()
}
