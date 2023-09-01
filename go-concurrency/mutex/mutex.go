package main

import (
	"fmt"
	"sync"
	"time"
)

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

func (c *SafeCounter) Inc(key string) {
	// When locked only one goroutine can access the counter at a time
	c.mu.Lock()
	c.v[key]++
	c.mu.Unlock()
}

func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v[key]
}

func main() {
	c := &SafeCounter{v: make(map[string]int)}
	for i := 0; i < 100000; i++ {
		go c.Inc("somekey")
	}

	fmt.Println(c.Value("somekey"))
	time.Sleep(time.Second)
	fmt.Println(c.Value("somekey"))
}
