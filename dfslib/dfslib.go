// IMPORTANT: This personal project implements the file system interface
// suggested by Ivan Beschastnikh @
// -- http://www.cs.ubc.ca/~bestchai/teaching/cs416_2017w2/assign2/index.html
// Unless otherwise noted, I certify all code is original and written by me

// Package dfslib specifies the interface by which applications
// access the distributed file system
package dfslib

import "fmt"

// Files are accessed in chunks of 32 bytes.
// Each file consists of 256 chunks.
type Chunk [32]byte
type chunkNum uint8

type FileMode int

// Files may be opened in any of the modes enumerated below
const (
	READ  FileMode = 1
	WRITE FileMode = 2
	// Disconnected read is currently not supported
	DREAD FileMode = 3
)

type DFSFile interface {
	Read(chunkNum uint8, chunk *Chunk) (err error)
	Write(chunkNum uint8, chunk *Chunk) (err error)
	Close() (err error)
}

type DFS interface {
	LocalFileExists(fname string) (exists bool, err error)
	GlobalFileExists(fname string) (exists bool, err error)
	Open(fname string, mode FileMode) (f DFSFile, err error)
	UMountDFS() (err error)
}

func MountDFS(serverAddr string, localIP string, localPath string) (dfs DFS, err error) {
	// TODO
	// For now return LocalPathError
	return nil, LocalPathError(localPath)
}

//==================================================================
// Error handling follows go conventions of explicitly typed errors.
// All errors returned by the DFS library are defined below.
//==================================================================

// The file path cannot be accessed or is otherwise invalid
type LocalPathError string

func (e LocalPathError) Error() string {
	return fmt.Sprintf("DFS: Cannot access local path [%s]", string(e))
}

// A user is identified by their IP and path
type UserRegistrationError string

func (e UserRegistrationError) Error() string {
	return fmt.Sprintf("server: The user: [%s] is already registered", string(e))
}

//==================================================================
// The client exposes an interface to the server. The server may
// invoke this interface to request information and data such as
// files stored locally on client or a specific chunk value.
// Client and server communicate via bi-directional RPC calls
//==================================================================

type ClientRPC int

type ClientInterface interface {
	Ping(stub int, reply *bool) (err error)
}
