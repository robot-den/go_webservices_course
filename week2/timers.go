package main

import (
	"fmt"
	"time"
)

func main() {
	//write_via_timer()
	
	write_via_ticker()
}

func write_via_timer() {
	timer := time.NewTimer(1 * time.Second)

	select {
	case <-timer.C:
		fmt.Println("timer.C ticked!")
	case <-time.After(3 * time.Second):
		fmt.Println("time.After timeout")
	case <-longSQLQuery():
		fmt.Println("Long query returned value to caller")
	}
}

func longSQLQuery() <- chan int {
	ch := make(chan int)
	go func(channel chan <- int ) {
		fmt.Println("Long query started")
		time.Sleep(5 * time.Second)
		ch <- 10
		fmt.Println("Long query finished")
	}(ch)
	return ch
}

func write_via_ticker()  {
	ticker := time.NewTicker(1 * time.Second)

	i := 0
	for tickTime := range ticker.C {
		if i >= 3 {
			ticker.Stop()
			break
		}
		fmt.Println("Ticker tick at time:", tickTime)
		i++
	}
	done := make(chan struct{})
	fmt.Println("Run AfterFunc")
	time.AfterFunc(2 * time.Second, longSQLQueryWrapperFunc(done))
	fmt.Println("Wait when func will finish")
	<-done
}

func longSQLQueryWrapperFunc(done chan<- struct{}) func() {
	return func() {
		result := longSQLQuery()
		fmt.Println("Waiting for long result")
		<-result
		close(done)
	}
}
