package main

import (
	"fmt"
	"time"
)

func ping(pings chan<- string, msg string) {
	pings <- msg
}

func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

func main_() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)

	ping(pings, "msg1")
	pong(pings, pongs)
	fmt.Println(<-pongs)

	c1 := make(chan string, 1)
	c2 := make(chan string, 1)

	go operation1(c1)
	go operation2(c2)
	close(c1)

	select {
	case msg1 := <-c1:
		fmt.Println(msg1)
	case msg2 := <-c2:
		fmt.Println(msg2)
	case <-time.Tick(1 * time.Second):
		fmt.Println("still waiting")
	}

}

func operation1(c1 chan string) {
	time.Sleep(time.Second * 2)
	c1 <- "one"
}

func operation2(c2 chan string) {
	time.Sleep(time.Second * 3)
	c2 <- "two"
}
