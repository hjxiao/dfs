/*
	Usage:
	go run server.go [server ip:port]
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
	// Disconnected read is currently not supported
	DREAD FileMode = 3
)

var (
	ipPort          string
	files           map[string]FileState             // Assumption: each file name is unique as namespace is global
	filesOpened     map[UserInfo]map[string]FileMode // Assumption: files cannot be deleted after opening
	registeredUsers []UserInfo
	rpcConnections  []net.Conn
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
	Unregister(stub int, reply *bool) (err error)
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
		fmt.Printf("server: Unable to bind to port [%s] to listen for incoming connection requests", ipPort)
		os.Exit(0)
	}

	for {
		conn, _ := listener.Accept()
		rpcConnections = append(rpcConnections, conn)
		fmt.Println("conns: ", rpcConnections)
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
	disconnectUser(user)
	fmt.Println("Users: ", registeredUsers) // TODO:
	fmt.Println("Conns: ", rpcConnections)  // TODO:
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
		*reply = true
		return nil
	}

	*reply = false
	return UserRegistrationError(user.LocalIP + " @ path " + user.LocalPath)
}

func (s *ServerRPC) Unregister(stub int, reply *bool) (err error) {
	return nil
}

func (s *ServerRPC) SendHeartbeat(user UserInfo, reply *bool) (err error) {
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

func disconnectUser(user UserInfo) {
	arrLen := len(rpcConnections)

	for i, conn := range rpcConnections {
		fmt.Println("Inside DisconnectUser")
		fmt.Println("conn addr: ", conn.RemoteAddr().String())
		fmt.Println("user LocalIP: ", user.LocalIP)
		if user.LocalIP == conn.LocalAddr().String() {
			conn.Close()

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
	return fmt.Sprintf("server: The user: [%s] is already registered", string(e))
}
