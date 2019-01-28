package learning

import (
	"context"
	"sync"
	"testing"
	"time"
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

func TestClosingAChannel(t *testing.T) {
	ch := make(chan int)
	// Publish 10 messages to the channel and then close the channel
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()
	for i := range ch {
		_ = i
	}
}

func TestChoosingBetweenChannelsWithRandomSend(t *testing.T) {
	chA := make(chan int, 1)
	chB := make(chan int, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		chA <- 1
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		chB <- 1
		wg.Done()
	}()
	select {
	case v := <-chA:
		_ = v // Got value from Channel A
	case v := <-chB:
		_ = v // Got value from Channel B
	}
	wg.Wait()
}

func TestChoosingBetweenChannelsToRead(t *testing.T) {
	chA := make(chan int, 1)
	chA <- 1
	chB := make(chan int, 1)
	chB <- 1
	select {
	case v := <-chA:
		_ = v // Got value from Channel A
	case v := <-chB:
		_ = v // Got value from Channel B
	}
}

func TestChoosingBetweenChannelsToReadWithDefault(t *testing.T) {
	chA := make(chan int)
	chB := make(chan int)
	gotDefault := false
	select {
	case v := <-chA:
		_ = v // Got value from Channel A
	case v := <-chB:
		_ = v // Got value from Channel B
	default:
		gotDefault = true
	}
	if !gotDefault {
		t.Fatal("Expected default to run but did not...")
	}
}

func TestReadingFromAClosedChannel(t *testing.T) {
	ch := make(chan int)
	close(ch)
	readFromClosedChannel := false
	select {
	case <-ch:
		readFromClosedChannel = true
	}
	if !readFromClosedChannel {
		t.Fatal("Reading from a closed channel should work and give a 0 value...")
	}
}

func TestCheckWhileReadingFromClosedChannel(t *testing.T) {
	ch := make(chan int)
	close(ch)
	_, ok := <-ch
	if ok {
		t.Error("Expected not OK but got OK...")
	}
}

func TestChoosingBetweenChannelsToWrite(t *testing.T) {
	chA := make(chan int)
	chB := make(chan int)
	readIntoChannel := ""
	// Only open up Channel A for reading
	go func() {
		<-chA
		readIntoChannel = "A"
	}()
	// Write to A or B channel - whichever is open
	select {
	case chA <- 1:
	case chB <- 1:
	}
	if readIntoChannel != "A" {
		t.Fatal("Expected channel A to have received a value...")
	}
}

func TestReadWriteOrTimeoutFromChannel(t *testing.T) {
	ch := make(chan int)
	timedOut := false
	select {
	case <-ch:
	case ch <- 1:
	case <-time.After(time.Millisecond):
		timedOut = true
	}
	if !timedOut {
		t.Fatal("Expected a timeout... didn't get it...")
	}
}

func TestCountElementsInChannel(t *testing.T) {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 1
	if len(ch) != 2 {
		t.Error("Expected 2 elements in channel but did not get.")
	}
	<-ch
	if len(ch) != 1 {
		t.Error("Expected 1 element1 in channel but did not get.")
	}
}

func TestContextCancellation(t *testing.T) {
	ctx := context.WithValue(context.Background(), paymentContextKey("confirmed"), make(chan struct{}))
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		time.Sleep(time.Millisecond * 100)
		cancel()
	}()
	status := ProcessPayment(ctx)
	if status != "CANCELLED" {
		t.Fatal("Expected a CANCELLED status.")
	}
}

func TestContextConfirmation(t *testing.T) {
	confirmationCh := make(chan struct{})
	ctx := context.WithValue(context.Background(), paymentContextKey("confirmed"), confirmationCh)
	go func() {
		time.Sleep(time.Millisecond * 100)
		confirmationCh <- struct{}{}
	}()
	status := ProcessPayment(ctx)
	if status != "CONFIRMED" {
		t.Fatal("Expected a CONFIRMED status.")
	}
}

type paymentContextKey string

// Refer http://blog.ralch.com/tutorial/golang-concurrency-patterns-context/
func ProcessPayment(ctx context.Context) string {
	confirmed := ctx.Value(paymentContextKey("confirmed")).(chan struct{})

	for {
		select {
		case <-confirmed:
			return "CONFIRMED"
		case <-ctx.Done():
			if ctx.Err() == context.Canceled {
				return "CANCELLED"
			}
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}
