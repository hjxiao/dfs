/*
	Usage:
	go run client.go [client IP:port] [server ip:port]
*/
package main

import (
	"fmt"
	"net/rpc"
	"os"
)

var (
	clientIPPort string
	serverIPPort string
)

func main() {
	args := os.Args[1:]
	clientIPPort = args[0]
	serverIPPort = args[1]

	client, err := rpc.Dial("tcp", serverIPPort)

	if err != nil {
		fmt.Println("client: Error establishing connection to server")
	}

	var reply bool
	fmt.Println("@@")
	err = client.Call("Server.Ping", 0, &reply)
	if err != nil {
		fmt.Println("client: Failed to ping server")
	}
}
