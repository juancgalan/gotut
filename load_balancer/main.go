package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

type DFAF func(*HTTPReq, byte) *HTTPReq

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

func sHTTPNewLineEOF(b byte, other DFAF) DFAF {
	if b == RETURN {
		return func(aa *HTTPReq, c byte) *HTTPReq {
			aa.message = append(aa.message, c)
			if c == NEWLINE {
				aa.eof = true
				aa.state = nil
				return aa
			} else {
				aa.state = other
				return aa
			}
		}

	} else {
		return other
	}
}

func sHTTPNewLine(b byte, newline, other DFAF) DFAF {
	if b == RETURN {
		return func(aa *HTTPReq, c byte) *HTTPReq {
			aa.message = append(aa.message, c)
			if c == NEWLINE {
				aa.state = newline
				return aa
			} else {
				aa.state = other
				return aa
			}
		}
	} else {
		return other
	}
}

func EOF(a *HTTPReq, b byte) *HTTPReq {
	a.message = append(a.message, b)
	a.state = sHTTPNewLineEOF(b, sHEADER)
	return a
}

func sHEADER(a *HTTPReq, b byte) *HTTPReq {
	a.message = append(a.message, b)
	a.state = sHTTPNewLine(b, EOF, sHEADER)
	return a
}

func sGET(a *HTTPReq, b byte) *HTTPReq {
	a.message = append(a.message, b)
	a.state = sHTTPNewLine(b, EOF, sGET)
	return a
}

func NewHTTPReq() (ans *HTTPReq) {
	ans = &HTTPReq{
		message: make([]byte, 0, 8192),
		state:   nil,
		eof:     false,
		err:     nil,
	}
	ans.state = sGET
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
		m = m.state(m, b)
		if m.err != nil {
			fmt.Printf("Invalid Answer\n")
			return nil, errors.New("500")
		}
		if m.eof {
			break
		}
	}
	for i := 0; i < 612; i += 1 {
		b, _ := conn.Rx.ReadByte()
		m.message = append(m.message, b)
	}
	return m.message, nil
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	m := NewHTTPReq()
	fmt.Println("--------Requested connection, acquiring data")
	for {
		b, _ := r.ReadByte()
		m = m.state(m, b)
		if m.err != nil {
			fmt.Printf("Invalid request\n")
			return
		}
		if m.eof {
			break
		}
	}
	fmt.Print(string(m.message))
	fmt.Println("--------Request received, rerouting")
	c, err := NewConnection("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Error redirecting to server")
		return
	}
	ans, _ := sendRequest(c, m.message)
	fmt.Print(string(ans))
	fmt.Println("--------Server responded, rerouting")
	w := bufio.NewWriter(conn)
	for i := range ans {
		w.WriteByte(ans[i])
	}
	w.Flush()
}

func SimpleBalancer() {
	ln, _ := net.Listen("tcp", ":8000")
	for {
		conn, _ := ln.Accept()
		handleRequest(conn)
	}
}

func main() {
	fmt.Printf("Starting Balancer...\n")
	fmt.Printf("^C for ending...\n")
	SimpleBalancer()
}
