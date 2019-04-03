package main

import (
	"bufio"
	"fmt"
	"net"
)

func SimpleBalancer() {
	ln, _ := net.Listen("tcp", ":80")

	conn, _ := ln.Accept()

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')

		fmt.Print("Message:", string(message))

	}
}

func main() {
	fmt.Printf("Starting Balancer...")
	SimpleBalancer()
}
