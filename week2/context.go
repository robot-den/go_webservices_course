package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ctx, finish := context.WithCancel(context.Background())
	resultCh := make(chan int, 1)
	fmt.Println("Run workers")
	for i := 0; i <= 10; i++  {
		go worker(ctx, i, resultCh)
	}
	fmt.Println("Wait for first result")
	result := <-resultCh
	fmt.Println("Result:", result)
	finish()

	fmt.Scanln()
}

func worker(ctx context.Context, i int, ch chan<-int) {
	waitTime := time.Duration(rand.Intn(100)+10) * time.Second
	fmt.Println("Worker", i, "sleep", waitTime)
	select {
	case <-ctx.Done():
		fmt.Println("Worker", i, "finished by context cancel")
	case <-time.After(waitTime):
		fmt.Println("Worker", i, "returns value")
		ch <- i
	}
}
