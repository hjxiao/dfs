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
	"sync"
	"time"
)

const (
	hbInterval = 5000 // defines heartbeat interval in milliseconds
)

type Chunk [32]byte

type FileMode int

// Files may be opened in any of the modes enumerated below
const (
	READ  FileMode = 1
	WRITE FileMode = 2
)

var (
	ipPort          string
	files           map[string]*FileState            // Assumption: global namespace, all file names are unique
	filesOpened     map[UserInfo]map[string]FileMode // Assumption: files cannot be deleted after opening
	clientConns     map[UserInfo]*rpc.Client
	registeredUsers []UserInfo
	lastHeartBeat   map[UserInfo]time.Time
)

type FileInfo struct {
	User  UserInfo
	Name  string
	Fmode FileMode
}

type FileState struct {
	fileExists       bool
	isLockedForWrite bool
	writeAccess      *sync.Mutex
	chunkVersion     []*FileVersionOwners // All chunks initialized at version 0; each write increments by 1
}

type FileVersionOwners struct {
	version int
	owners  []UserInfo
}

type UserInfo struct {
	LocalIP   string
	LocalPath string
}

type WriteInfo struct {
	User     UserInfo
	Fname    string
	ChunkNum uint8
}

type ReadInfo struct {
	User          UserInfo
	Fname         string
	ChunkNum      uint8
	LocalChunkVer int
}

type ReadValue struct {
	Chnk           Chunk
	IsNew          bool
	globalChunkVer int
}

type ServerRPC int

type ServerInterface interface {
	Ping(stub int, reply *bool) (err error)
	Register(user UserInfo, reply *bool) (err error)
	Unregister(user UserInfo, reply *bool) (err error)
	SendHeartbeat(user UserInfo, reply *bool) (err error)
	EstablishReverseRPC(user UserInfo, reply *bool) (err error)
	FileExists(fname string, reply *bool) (err error)
	RegisterFile(fi FileInfo, reply *bool) (err error)
	WriteFile(wi WriteInfo, reply *bool) (err error)
	ReadFile(ri ReadInfo, rv *ReadValue) (err error)
	CloseFile(fi FileInfo, reply *bool) (err error)
}

func main() {
	args := os.Args[1:]
	fmt.Println("args: ", args)
	ipPort = args[0]

	files = make(map[string]*FileState, 0)
	filesOpened = make(map[UserInfo]map[string]FileMode, 0)
	clientConns = make(map[UserInfo]*rpc.Client, 0)
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

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
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

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
func reap(user UserInfo) {
	fmt.Printf("server: [%s] disconnected due to late heartbeat\n", user)
	removeUser(user)
	fmt.Println("Users: ", registeredUsers)
}

//==================================================================
// Server interface
//==================================================================

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
func (s *ServerRPC) Ping(stub int, reply *bool) (err error) {
	fmt.Println("Received ping from client")
	return nil
}

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
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

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
func (s *ServerRPC) Unregister(user UserInfo, reply *bool) (err error) {
	fmt.Printf("server: Removing requested user [%s]\n", user)
	removeUser(user)
	fmt.Println("Users: ", registeredUsers)
	return nil
}

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
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

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
func (s *ServerRPC) EstablishReverseRPC(user UserInfo, reply *bool) (err error) {
	clientConns[user], err = rpc.Dial("tcp", user.LocalIP)
	if err != nil {
		return err
	}

	r := false
	clientConns[user].Call("ClientRPC.Ping", 0, &r)

	return nil
}

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
func (s *ServerRPC) FileExists(fname string, reply *bool) (err error) {
	fs := files[fname]
	if fs == nil {
		*reply = false
	} else {
		fmt.Println("server: File existence - ", fs.fileExists)
		*reply = fs.fileExists
	}
	return nil
}

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
func (s *ServerRPC) RegisterFile(fi FileInfo, reply *bool) (err error) {
	createFileIfNotExist(fi)

	err = configureWriteAccess(fi)
	if err != nil {
		return err
	}

	updateOpenedFiles(fi)
	return nil
}

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
func (s *ServerRPC) WriteFile(wi WriteInfo, reply *bool) (err error) {
	fs := files[wi.Fname]
	cv := fs.chunkVersion
	fvo := cv[wi.ChunkNum]

	if fvo == nil {
		fvo = &FileVersionOwners{version: 0, owners: make([]UserInfo, 0)}
		cv[wi.ChunkNum] = fvo
	}

	fvo.version++
	fvo.owners = make([]UserInfo, 0)
	fvo.owners = append(fvo.owners, wi.User)

	*reply = true
	return nil
}

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
func (s *ServerRPC) ReadFile(ri ReadInfo, rv *ReadValue) (err error) {
	fs := files[ri.Fname]
	cv := fs.chunkVersion
	fvo := cv[ri.ChunkNum]

	if fvo == nil {
		fvo = &FileVersionOwners{version: 0, owners: make([]UserInfo, 0)}
		cv[ri.ChunkNum] = fvo
	}

	if ri.LocalChunkVer < fvo.version {
		for _, user := range fvo.owners {
			readInfoForRetrievingChunk := ReadInfo{User: user, Fname: ri.Fname, ChunkNum: ri.ChunkNum}
			newChunk, err := retrieveLatestChunk(readInfoForRetrievingChunk)
			if err != nil {
				continue
			} else {
				rv.Chnk = newChunk
				rv.globalChunkVer = fvo.version
				rv.IsNew = true
				break
			}
		}
	} else {
		rv.IsNew = false
	}

	if !containsUser(ri.User, fvo.owners) {
		fvo.owners = append(fvo.owners, ri.User)
	}

	return nil
}

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
func (s *ServerRPC) CloseFile(fi FileInfo, reply *bool) (err error) {
	if fi.Fmode == WRITE {
		fs := files[fi.Name]
		fs.isLockedForWrite = false
		fs.writeAccess.Unlock()
	}
	*reply = true
	return nil
}

//==================================================================
// Helper Functions
//==================================================================

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
func containsUser(user UserInfo, regUsers []UserInfo) bool {
	for _, regUser := range regUsers {
		if userEquals(user, regUser) {
			return true
		}
	}
	return false
}

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
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

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
func userEquals(u, ru UserInfo) bool {
	return (u.LocalIP == ru.LocalIP) && (u.LocalPath == ru.LocalPath)
}

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
func createFileIfNotExist(fi FileInfo) {
	if files[fi.Name] == nil {
		fs := FileState{fileExists: true,
			isLockedForWrite: false,
			writeAccess:      &sync.Mutex{},
			chunkVersion:     make([]*FileVersionOwners, 256)}

		files[fi.Name] = &fs
	}
}

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
func configureWriteAccess(fi FileInfo) error {
	if fi.Fmode == WRITE {
		if files[fi.Name].isLockedForWrite == true {
			return OpenWriteConflictError(fi.Name)
		}
		files[fi.Name].isLockedForWrite = true
		files[fi.Name].writeAccess.Lock()
	}

	return nil
}

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
func updateOpenedFiles(fi FileInfo) {
	if filesOpened[fi.User] == nil {
		filesOpened[fi.User] = make(map[string]FileMode, 0)
	}

	filesOpened[fi.User][fi.Name] = fi.Fmode
}

/*
 Purpose:
 Params:
 Returns
 Throws:
*/
func retrieveLatestChunk(ri ReadInfo) (c Chunk, err error) {
	connToClient := clientConns[ri.User]

	if containsUser(ri.User, registeredUsers) && connToClient != nil {
		err = connToClient.Call("ClientRPC.RetrieveLatestChunk", ri, &c)
		if err != nil {
			return c, ChunkUnavailableError(ri.ChunkNum)
		}
	} else {
		return c, ChunkUnavailableError(ri.ChunkNum)
	}

	return c, nil
}

//==================================================================
// Errors
//==================================================================

// Contains chunkNum that is unavailable
type ChunkUnavailableError uint8

func (e ChunkUnavailableError) Error() string {
	return fmt.Sprintf("DFS: Latest verson of chunk [%d] unavailable", e)
}

// A user is identified by their IP, port, and file path
type UserRegistrationError string

func (e UserRegistrationError) Error() string {
	return fmt.Sprintf("server: The user: [%s] is already registered\n", string(e))
}

type HeartbeatRegistrationError string

func (e HeartbeatRegistrationError) Error() string {
	return fmt.Sprintf("server: The user [%s] sent a heartbeat, but is not registered\n", string(e))
}

// Contains filename
type OpenWriteConflictError string

func (e OpenWriteConflictError) Error() string {
	return fmt.Sprintf("DFS: Filename [%s] is opened for writing by another client", string(e))
}
