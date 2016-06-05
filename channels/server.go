package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

var randCompute *rand.Rand

type ServerHandler struct {
	listener      net.Listener
	server        *http.Server
	counter       uint64
	shutdownState uint32
	counterDone   chan bool
}

type ClosableListener struct {
	net.Listener
	close chan bool
}

func NewServerHandler(address string) (*ServerHandler, error) {
	sh := &ServerHandler{}
	sh.server = &http.Server{
		Handler: sh,
		ConnState: func(conn net.Conn, state http.ConnState) {
			switch state {
			case http.StateNew:
				atomic.AddUint64(&sh.counter, uint64(1))
			case http.StateClosed:
				atomic.AddUint64(&sh.counter, ^uint64(0))
			}
			if atomic.LoadUint32(&sh.shutdownState) == 1 && atomic.LoadUint64(&sh.counter) == 0 {
				sh.counterDone <- true
			}

		},
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	sh.listener = listener
	sh.counter = 0
	sh.shutdownState = 0
	sh.counterDone = make(chan bool)

	return sh, nil
}

func init() {
	randCompute = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
}

func main() {
	sh, err := NewServerHandler(":5555")
	if err != nil {
		log.Fatal(err)
	}

	sh.server.Serve(sh.listener)
	select {
	case <-sh.counterDone:
		break
	}
}

func handleRequest(r *http.Request) {
	fmt.Printf("I started handling %p\n", r)
	_ = randCompute.Intn(3000)
	time.Sleep(10 * time.Second)
	fmt.Printf("Done processing %p\n", r)

}

func (sh *ServerHandler) shutdown() {
	sh.listener.Close()
	sh.server.SetKeepAlivesEnabled(false)
	atomic.StoreUint32(&sh.shutdownState, uint32(1))

	return
}

func (sh *ServerHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/shutdown" {
		fmt.Println("Shut down accepted")

		sh.shutdown()
		return
	}

	handleRequest(r)
}
