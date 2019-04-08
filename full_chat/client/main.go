package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/chzyer/readline"
)

type Connection struct {
	Conn net.Conn
	Tx   *bufio.Writer
	Rx   *bufio.Reader
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

func ClearConnection(con *Connection) {
	if con == nil {
		return
	}
	con.Conn.Close()
	con = nil
}

func send(con *Connection, message string) {
	if con == nil {
		fmt.Println("Error, connection not active.")
	}
	fmt.Println("Sending message... ", message)
	b := []byte(message)
	for i := range b {
		con.Tx.WriteByte(b[i])
	}
	con.Tx.Flush()
}

func listen(con *Connection) {
	for {
		b, err := con.Rx.ReadByte()
		if err != nil {
			fmt.Println("Error reading client: ", err.Error())
			return
		}
		fmt.Printf("%c", b)
	}
}

func connect(ctx string) *Connection {
	s := strings.Split(ctx, " ")
	var prot, addr string
	switch {
	case len(s) == 1:
		prot = "tcp"
		addr = s[0]
		break
	case len(s) == 2:
		prot = s[0]
		addr = s[1]
		break
	default:
		fmt.Println("Error connection format. Example: conn [protocol] address[:port]")
	}
	ncon, _ := NewConnection(prot, addr)
	if ncon != nil {
		go listen(ncon)
	}
	return ncon

}

func main() {
	l, err := readline.NewEx(&readline.Config{
		Prompt:              ">>",
		HistoryFile:         "read.tmp",
		AutoComplete:        nil,
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   false,
		FuncFilterInputRune: nil,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	read := true
	var connection *Connection
	for read {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		}
		switch {
		case strings.HasPrefix(line, "conn"):
			fmt.Println("%[", line, "]")
			line = strings.TrimSpace(line)
			connection = connect(strings.TrimSpace(strings.TrimPrefix(line, "conn")))
			break
		case strings.HasPrefix(line, "exit"):
			fmt.Println("%[", line, "]")
			line = strings.TrimSpace(line)
			fmt.Println("Bye")
			read = false
			break
		default:
			send(connection, line)
			break
		}
	}
}
