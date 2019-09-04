package learning

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func infiniteJob(ctx context.Context, i int) {
	for {
		select {
		case <-time.After(time.Second):
			log.Printf("job %d is alive", i)
		case <-ctx.Done():
			return
		}
	}
}

func TestCtx(t *testing.T) {
	if os.Getenv("TEST_CTX") == "" {
		t.Skip("Set the TEST_CTX flag to run this test.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			infiniteJob(ctx, i)
		}(i)
	}

	wg.Wait()
}
