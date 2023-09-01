package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

func Walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}

	Walk(t.Left, ch)
	ch <- t.Value
	Walk(t.Right, ch)
}

func Walking(t *tree.Tree, ch chan int) {
	defer close(ch)
	Walk(t, ch)
}

func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go Walking(t1, ch1)
	go Walking(t2, ch2)

	for {
		v1, ok1 := <- ch1
		v2, ok2 := <- ch2
		if ok1 != ok2 || v1 != v2 {
			return false
		}
		if !ok1 {
			break
		}
	}

	return true
}

func main() {
	t1 := tree.New(1)
	t2 := tree.New(1)
	same := Same(t1, t2)
	fmt.Println(same)
}
