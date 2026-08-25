package middleware

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRateLimiterConcurrentWindowSweep(t *testing.T) {
	la := NewRateLimiter(200, time.Minute)
	lb := NewRateLimiter(200, time.Minute)
	base := time.Now().UTC()
	la.mu.Lock()
	la.lastSweep = base.Add(-2 * time.Minute)
	la.mu.Unlock()

	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			l := la
			if i%2 == 1 {
				l = lb
			}
			for j := 0; j < 300; j++ {
				l.Allow(fmt.Sprintf("key-%d", i%4), base.Add(time.Duration(j%40)*time.Second))
			}
		}(i)
	}
	close(start)
	wg.Wait()
}

func TestRateLimiterIndependentMaps(t *testing.T) {
	la := NewRateLimiter(2, time.Hour)
	lb := NewRateLimiter(2, time.Hour)
	now := time.Now().UTC()
	for i := 0; i < 2; i++ {
		if ok, _ := la.Allow("shared-key", now); !ok {
			t.Fatalf("la request %d should be allowed", i+1)
		}
	}
	if ok, _ := lb.Allow("shared-key", now); !ok {
		t.Fatal("lb should keep an independent counter")
	}
	if ok, _ := lb.Allow("shared-key", now); !ok {
		t.Fatal("lb second request should be allowed")
	}
	if ok, _ := lb.Allow("shared-key", now); ok {
		t.Fatal("lb third request should be rate-limited")
	}
}
