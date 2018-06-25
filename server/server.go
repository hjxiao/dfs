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

const (
	hbInterval = 2000 // defines heartbeat interval in milliseconds
)

type FileMode int

// Files may be opened in any of the modes enumerated below
const (
	READ  FileMode = 1
	WRITE FileMode = 2
	// Disconnected read is currently not supported
	DREAD FileMode = 3
)

var (
	ipPort          string
	files           map[string]FileState // Assumption: each file name is unique as namespace is global
	filesOpened     map[UserInfo]map[string]FileMode
	registeredUsers []UserInfo
)

type FileState struct {
	fileExists   bool
	chunkVersion [256]FileVersionOwners // All chunks initialized at version 0; each write increments by 1
}

type FileVersionOwners struct {
	version int
	owners  []UserInfo
}

type UserInfo struct {
	LocalIP   string
	LocalPath string
}

type ServerRPC int

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

	server := new(ServerRPC)
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

//==================================================================
// Monitor reaps registered users who have either
// (1) failed to send a heartbeat or
// (2) sent a late heartbeat
//==================================================================

func monitor() {

}

//==================================================================
// Server interface
//==================================================================

func (s *ServerRPC) Ping(stub int, reply *bool) (err error) {
	fmt.Println("Received ping from client")
	return nil
}

func (s *ServerRPC) Register(user UserInfo, reply *bool) (err error) {
	if !containsUser(user, registeredUsers) {
		registeredUsers = append(registeredUsers, user)
		*reply = true
		fmt.Println("@@")
		fmt.Println(registeredUsers)
		fmt.Println("@@")
		return nil
	}

	*reply = false
	return UserRegistrationError(user.LocalIP + " & " + user.LocalPath)
}

func (s *ServerRPC) Unregister(stub int, reply *bool) (err error) {
	return nil
}

func (s *ServerRPC) SendHeartbeat(stub int, reply *bool) (err error) {
	return nil
}

//==================================================================
// Helper Functions
//==================================================================

func containsUser(user UserInfo, regUsers []UserInfo) bool {
	for _, regUser := range regUsers {
		if userEquals(user, regUser) {
			return true
		}
	}
	return false
}

func removeUser(user UserInfo, regUsers []UserInfo) {
	for i, regUser := range regUsers {
		if userEquals(user, regUser) {
			regUsers[i] = regUsers[len(regUsers)-1]
			regUsers = regUsers[:len(regUsers)-1]
			break
		}
	}
}

func userEquals(u, ru UserInfo) bool {
	return (u.LocalIP == ru.LocalIP) && (u.LocalPath == ru.LocalPath)
}

//==================================================================
// Errors
//==================================================================

// A user is identified by their IP, port, and file path
type UserRegistrationError string

func (e UserRegistrationError) Error() string {
	return fmt.Sprintf("server: The user: [%s] is already registered", string(e))
}
