package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func requestHandler(conn net.Conn) {
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message:", string(message))
}

func main() {

	hostPtr := flag.String("h", "localhost", "Remote Address")
	portPtr := flag.String("p", "80", "Remote Port")
	flag.Parse()
	l, err := net.Listen("tcp", *hostPtr+":"+*portPtr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + *hostPtr + ":" + *portPtr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Unexpected Error", err.Error())
			os.Exit(1)
		}
		go requestHandler(conn)
	}
}
