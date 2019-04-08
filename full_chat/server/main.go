package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type ObservableByte struct {
	data      byte
	observers []chan byte
	lock      sync.Mutex
}

func (o *ObservableByte) subscribe(c *chan byte) {
	o.observers = append(o.observers, *c)
}

func (o *ObservableByte) set(newData byte) {
	o.lock.Lock()
	o.data = newData
	for i := range o.observers {
		o.observers[i] <- newData
	}
	o.lock.Unlock()
}

func NewObservableByte() ObservableByte {
	return ObservableByte{
		data:      0,
		observers: make([]chan byte, 0),
		lock:      sync.Mutex{},
	}
}

func ClientTxHandler(conn net.Conn, out chan byte) {
	w := bufio.NewWriter(conn)
	for {
		b := <-out
		w.WriteByte(b)
		w.Flush()
	}
}

func ClientRxHandler(conn net.Conn, buffer *ObservableByte) {
	r := bufio.NewReader(conn)
	for {
		b, err := r.ReadByte()
		if err != nil {
			fmt.Println("Error reading client: ", err.Error())
			return
		}
		buffer.set(b)
	}
}

func (o *ObservableByte) handle(conn net.Conn) {
	c := make(chan byte, 8192)
	o.subscribe(&c)
	fmt.Println("Client connected")
	go ClientTxHandler(conn, c)
	go ClientRxHandler(conn, o)
	fmt.Fprintf(conn, "Welcome...\n")
}

func LocalEcho(buffer *ObservableByte) {
	c := make(chan byte)
	buffer.subscribe(&c)
	for {
		i := <-c
		fmt.Printf("%c", i)
	}
}

func SimpleEchoServer() {
	ln, _ := net.Listen("tcp", ":80")
	ob := NewObservableByte()
	go LocalEcho(&ob)
	for {
		conn, _ := ln.Accept()
		ob.handle(conn)
	}
}

func main() {
	fmt.Printf("Starting Chat Server...\n")
	fmt.Printf("^C for ending...")
	SimpleEchoServer()
}
