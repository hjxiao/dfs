/*
	Usage:
	go run server.go [server ip:port]
*/
package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
)

var (
	ipPort = "" // port to listen for incoming client connection requests
)

type UserInfo struct {
	localIP   string
	localPath string
}

type Server int

type ServerInterface interface {
	Ping(stub int, reply *bool) (err error)
	Register(user UserInfo, reply *bool) (err error)
	Unregister(stub int, reply *bool) (err error)
	SendHeartbeat(stub int, reply *bool) (err error)
}

func main() {
	args := os.Args[1:]
	fmt.Println("args: ", args)
	ipPort = args[0]

	server := new(Server)
	serverRPC := rpc.NewServer()
	serverRPC.Register(server)

	listener, err := net.Listen("tcp", ipPort)
	if err != nil {
		fmt.Printf("server: Unable to bind to port [%s] to listen for incoming connection requests", ipPort)
		os.Exit(0)
	}

	for {
		conn, _ := listener.Accept()
		go serverRPC.ServeConn(conn)
	}

}

func (s *Server) Ping(stub int, reply *bool) (err error) {
	fmt.Println("Received ping from client")
	return nil
}

func (s *Server) Register(user UserInfo, reply *bool) (err error) {
	return nil
}

func (s *Server) Unregister(stub int, reply *bool) (err error) {
	return nil
}

func (s *Server) SendHeartbeat(stub int, reply *bool) (err error) {
	return nil
}
