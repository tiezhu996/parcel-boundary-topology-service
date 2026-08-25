package geometry

import (
	"sync"
	"testing"
)

func TestParsePolygonConcurrentCache(t *testing.T) {
	raw := polygonJSON(`[0,0],[10,0],[10,10],[0,10],[0,0]`)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 200; j++ {
				if _, err := ParsePolygon(raw); err != nil {
					t.Errorf("ParsePolygon() error = %v", err)
					return
				}
			}
		}()
	}
	close(start)
	wg.Wait()
}
