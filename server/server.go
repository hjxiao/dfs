/*
	Usage:
	go run server.go [server ip:port]

	Example:
	go run server.go 127.0.0.1:3000
*/
package main

import (
	"fmt"
	"net"
	rpc "net/rpc"
	"os"
	"time"
)

const (
	hbInterval = 5000 // defines heartbeat interval in milliseconds
)

type FileMode int

// Files may be opened in any of the modes enumerated below
const (
	READ  FileMode = 1
	WRITE FileMode = 2
)

var (
	ipPort          string
	files           map[string]FileState             // Assumption: each file name is unique as namespace is global
	filesOpened     map[UserInfo]map[string]FileMode // Assumption: files cannot be deleted after opening
	registeredUsers []UserInfo
	lastHeartBeat   map[UserInfo]time.Time
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
	Unregister(user UserInfo, reply *bool) (err error)
	SendHeartbeat(user UserInfo, reply *bool) (err error)
}

func main() {
	args := os.Args[1:]
	fmt.Println("args: ", args)
	ipPort = args[0]

	lastHeartBeat = make(map[UserInfo]time.Time, 0)

	server := new(ServerRPC)
	serverRPC := rpc.NewServer()
	serverRPC.Register(server)

	listener, err := net.Listen("tcp", ipPort)
	if err != nil {
		fmt.Printf("server: Unable to bind to port [%s] to listen for incoming connection requests\n", ipPort)
		os.Exit(0)
	}

	for {
		conn, _ := listener.Accept()
		go serverRPC.ServeConn(conn)
	}
}

//==================================================================
// Monitor reaps a user if it
// (1) failed to send a heartbeat or
// (2) sent a late heartbeat (> 2 seconds)
//==================================================================

func monitor(user UserInfo) {
	for {
		timeBetween := time.Now().Sub(lastHeartBeat[user])
		if timeBetween > hbInterval*time.Millisecond {
			reap(user)
			break
		} else {
			time.Sleep(2 * time.Second)
		}
	}
}

func reap(user UserInfo) {
	fmt.Printf("server: [%s] disconnected due to late heartbeat\n", user)
	removeUser(user)
	fmt.Println("Users: ", registeredUsers)
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
		lastHeartBeat[user] = time.Now()
		go monitor(user)
		fmt.Println("server: Received register from ", user)
		*reply = true
		return nil
	}

	*reply = false
	return UserRegistrationError(user.LocalIP + " @ path " + user.LocalPath)
}

func (s *ServerRPC) Unregister(user UserInfo, reply *bool) (err error) {
	fmt.Printf("server: Removing requested user [%s]\n", user)
	removeUser(user)
	fmt.Println("Users: ", registeredUsers)
	return nil
}

func (s *ServerRPC) SendHeartbeat(user UserInfo, reply *bool) (err error) {
	if !containsUser(user, registeredUsers) {
		*reply = false
		return HeartbeatRegistrationError(user.LocalIP + " @ path " + user.LocalPath)
	}

	fmt.Println("server: Received heartbeat from: ", user)
	lastHeartBeat[user] = time.Now()
	*reply = true
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

func removeUser(user UserInfo) {
	arrLen := len(registeredUsers)

	for i, regUser := range registeredUsers {
		if userEquals(user, regUser) {
			if arrLen == 1 {
				registeredUsers = make([]UserInfo, 0)
				break
			}

			registeredUsers[i] = registeredUsers[arrLen-1]
			registeredUsers = registeredUsers[:arrLen-1]
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
	return fmt.Sprintf("server: The user: [%s] is already registered\n", string(e))
}

type HeartbeatRegistrationError string

func (e HeartbeatRegistrationError) Error() string {
	return fmt.Sprintf("server: The user [%s] sent a heartbeat, but is not registered\n", string(e))
}
