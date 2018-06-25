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
	path         string = "/tmp/"
)

type UserInfo struct {
	LocalIP   string
	LocalPath string
}

func main() {
	args := os.Args[1:]
	clientIPPort = args[0]
	serverIPPort = args[1]

	myUser := UserInfo{LocalIP: clientIPPort, LocalPath: path}

	client, err := rpc.Dial("tcp", serverIPPort)

	if err != nil {
		fmt.Println("client: Error establishing connection to server")
	}

	var reply bool
	fmt.Println("@@")
	err = client.Call("ServerRPC.Ping", 0, &reply)
	if err != nil {
		fmt.Println("client: Failed to ping server")
	}

	var reply2 bool
	err = client.Call("ServerRPC.Register", myUser, &reply2)
	if err != nil {
		fmt.Println("client: Failed to register with server")
		fmt.Println(err.Error())
	}
}
