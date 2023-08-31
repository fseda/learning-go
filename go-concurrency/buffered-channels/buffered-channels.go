package main

import "fmt"

func fibonacci(n int, c chan uint64) {
	var x, y uint64

	x, y = 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}

func main() {
	ch1 := make(chan int, 2)
	ch1 <- 1
	ch1 <- 2

	ch2 := make(chan uint64, 10)
	go fibonacci(cap(ch2), ch2)
	for i := range ch2 {
		fmt.Println(i)
	}
}
