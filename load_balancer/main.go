package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

type DFAF func(byte) *HTTPReq

const NEWLINE = 10
const RETURN = 13

type Connection struct {
	Conn net.Conn
	Tx   *bufio.Writer
	Rx   *bufio.Reader
}

type HTTPReq struct {
	message []byte
	state   DFAF
	eof     bool
	err     error
}

func (a *HTTPReq) waitEOF(b byte) *HTTPReq {
	a.message = append(a.message, b)
	if b == NEWLINE {
		a.eof = true
		a.state = nil
	} else {
		a.state = a.reqBody
	}
	return a
}

func (a *HTTPReq) endBody(b byte) *HTTPReq {
	a.message = append(a.message, b)
	if b == RETURN {
		a.state = a.waitEOF
	} else {
		a.state = a.reqBody
	}
	return a
}

func (a *HTTPReq) waitNewLine(b byte) *HTTPReq {
	a.message = append(a.message, b)
	if b == NEWLINE {
		a.state = a.endBody
	} else {
		a.state = a.reqBody
	}
	return a
}

func (a *HTTPReq) reqBody(b byte) *HTTPReq {
	a.message = append(a.message, b)
	if b == RETURN {
		a.state = a.waitNewLine
	} else {
		a.state = a.reqBody
	}
	return a
}

func (a *HTTPReq) init(b byte) *HTTPReq {
	if b == RETURN {
		a.state = nil
		a.err = errors.New("Invalid Request!")
	}
	a.message = append(a.message, b)
	a.state = a.reqBody
	return a
}

func NewHTTPReq() (ans *HTTPReq) {
	ans = &HTTPReq{
		message: make([]byte, 0, 8192),
		state:   nil,
		eof:     false,
		err:     nil,
	}
	ans.state = ans.init
	return
}

func NewConnection(protocol, address string) (*Connection, error) {
	conn, err := net.Dial(protocol, address)
	if err != nil {
		fmt.Println("Error connecting...", err.Error())
		return nil, err
	}
	return &Connection{
		Conn: conn,
		Tx:   bufio.NewWriter(conn),
		Rx:   bufio.NewReader(conn),
	}, nil
}

func sendRequest(conn *Connection, data []byte) ([]byte, error) {
	for i := range data {
		conn.Tx.WriteByte(data[i])
	}
	conn.Tx.Flush()
	m := NewHTTPReq()
	for {
		b, _ := conn.Rx.ReadByte()
		m = m.state(b)
		if m.err != nil {
			fmt.Printf("Invalid Answer\n")
			return nil, errors.New("500")
		}
		if m.eof {
			break
		}
	}
	return m.message, nil
}

func handleRequest(conn net.Conn) {
	r := bufio.NewReader(conn)
	m := NewHTTPReq()
	for {
		b, _ := r.ReadByte()
		m = m.state(b)
		if m.err != nil {
			fmt.Printf("Invalid request\n")
			return
		}
		if m.eof {
			break
		}
	}
	fmt.Println("Requested received, rerouting")
	c, err := NewConnection("tcp", "localhost:80")
	if err != nil {
		fmt.Printf("Error redirecting to server")
		return
	}
	ans, _ := sendRequest(c, m.message)
	fmt.Println("Server responded, rerouting")
	fmt.Print(string(ans))
}

func SimpleBalancer() {
	ln, _ := net.Listen("tcp", ":8080")
	for {
		conn, _ := ln.Accept()
		handleRequest(conn)
	}
}

func main() {
	fmt.Printf("Starting Balancer...\n")
	fmt.Printf("^C for ending...")
	SimpleBalancer()
}
