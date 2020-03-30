package main

import (
	"fmt"
)

func main() {
	//closeAfterFirstMessage()
	readFromChannelWhileItOpen()
}

func closeAfterFirstMessage() {
	ch := make(chan struct{})

	go func(channel <-chan struct{}) {
		fmt.Println("Start separate gorutine")
		<- channel
		fmt.Println("Stop separate gorutine")
	}(ch)

	fmt.Scanln()
	fmt.Println("Write to channel from main process")
	ch <- struct{}{}
	fmt.Scanln()
}

func readFromChannelWhileItOpen()  {
	ch := make(chan string)
	done := make(chan struct{}, 5)

	go func(channel chan<- string, d <-chan struct{}) {
		defer close(channel)

		fmt.Println("Start separate gorutine")
		for i := 0; i <= 5 ; i++  {
			msg := fmt.Sprintf("v: %v", i)

			select {
			case channel <- msg:
				fmt.Println("  Write to channel:", msg)
			case <-d:
				fmt.Println("Stop separate gorutine by command")
				return
			}
		}
		fmt.Println("Stop separate gorutine by itself")
	}(ch, done)

	fmt.Scanln()
	fmt.Println("Read from channel in main process")
	for value := range ch {
		fmt.Println("  Read from channel:", value)
	}

	close(done)

	fmt.Scanln()
}
