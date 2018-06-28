/*
	Usage:
	go run client.go [client IP:port] [server ip:port]

	Example:
	go run client.go 127.0.0.1:3001 127.0.0.1:3000
*/
package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"time"
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
	fmt.Println(net.InterfaceAddrs())

	if err != nil {
		fmt.Println("client: Error establishing connection to server")
	}

	var reply bool
	err = client.Call("ServerRPC.Ping", 0, &reply)
	if err != nil {
		fmt.Println("client: Failed to ping server")
	}

	var reply2 bool
	err = client.Call("ServerRPC.Register", myUser, &reply2)
	// err = client.Call("ServerRPC.Register", myUser, &reply2)
	// if err != nil {
	// 	fmt.Println("client: Failed to register with server")
	// 	fmt.Println(err.Error())
	// }
	for {
		var reply3 bool
		fmt.Println("@@")
		err = client.Call("ServerRPC.SendHeartbeat", myUser, &reply3)
		fmt.Println("##")
		if err != nil {
			fmt.Println(err.Error()) // failure detector on client side
			break
		}
		time.Sleep(time.Millisecond * 2500)
		err = client.Call("ServerRPC.Unregister", myUser, &reply3)
	}
}
