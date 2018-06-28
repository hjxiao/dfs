// IMPORTANT: This personal project implements the file system interface
// suggested by Ivan Beschastnikh @
// -- http://www.cs.ubc.ca/~bestchai/teaching/cs416_2017w2/assign2/index.html
// Unless otherwise noted, I certify all code is original and written by me

// Package dfslib specifies the interface by which applications
// access the distributed file system
package dfslib

import (
	"fmt"
	"net"
	"os"
)

// Files are accessed in chunks of 32 bytes.
// Each file consists of 256 chunks.
type Chunk [32]byte

type FileMode int

// Files may be opened in any of the modes enumerated below
const (
	READ  FileMode = 1
	WRITE FileMode = 2
	// Disconnected read is currently not supported
	DREAD FileMode = 3
)

var (
	theDFSInstance DFS // singleton pattern
	connToServer   *net.Conn
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

type dfsObject struct {
}

func MountDFS(serverAddr string, localIP string, localPath string) (dfs DFS, err error) {
	if checkLocalPathOK(localPath) {
		if theDFSInstance != nil {
			theDFSInstance = dfsObject{}
		}
		return theDFSInstance, nil
	}
	return nil, LocalPathError(localPath)
}

//================================
// MountDFS Helper Functions
//================================

func checkLocalPathOK(localPath string) bool {
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func connectToServer(sAddr string) error {
	client, err := net.Dial("tcp", sAddr)
	if err != nil {
		return err
	}

	connToServer = &client
	return nil
}

//================================
// IMPLEMENTATION: DFS interface
//================================

func (dfs dfsObject) LocalFileExists(fname string) (exists bool, err error) {
	// TODO:
	return false, NotImplementedError("DFS.LocalFileExists")
}

func (dfs dfsObject) GlobalFileExists(fname string) (exists bool, err error) {
	// TODO:
	return false, NotImplementedError("DFS.GlobalFileExists")
}

func (dfs dfsObject) Open(fname string, mode FileMode) (f DFSFile, err error) {
	// TODO:
	return nil, NotImplementedError("DFS.UMountDFS")
}

func (dfs dfsObject) UMountDFS() (err error) {
	// TODO:
	return NotImplementedError("DFS.UMountDFS")
}

//===================================
// IMPLEMENTATION: DFSFile Interface
//===================================

func Read(chunkNum uint8, chunk *Chunk) (err error) {
	// TODO:
	return NotImplementedError("DFSFile.Read")
}

func Write(chunkNum uint8, chunk *Chunk) (err error) {
	// TODO:
	return NotImplementedError("DFSFile.Write")
}

func Close() (err error) {
	// TODO:
	return NotImplementedError("DFSFile.Close")
}

//==================================================================
// Error handling follows go conventions of explicitly typed errors.
// All errors returned by the DFS library are defined below.
//==================================================================

// A user is identified by their IP and path
type UserRegistrationError string

func (e UserRegistrationError) Error() string {
	return fmt.Sprintf("dfs: The user: [%s] is already registered on server", string(e))
}

// Contains chunkNum that is unavailable
type ChunkUnavailableError uint8

func (e ChunkUnavailableError) Error() string {
	return fmt.Sprintf("DFS: Latest verson of chunk [%d] unavailable", e)
}

// Contains filename
type OpenWriteConflictError string

func (e OpenWriteConflictError) Error() string {
	return fmt.Sprintf("DFS: Filename [%s] is opened for writing by another client", string(e))
}

// Contains file mode that is bad.
type BadFileModeError FileMode

func (e BadFileModeError) Error() string {
	return fmt.Sprintf("DFS: Cannot perform this operation in current file mode [%s]", string(e))
}

// Contains filename.
type WriteModeTimeoutError string

func (e WriteModeTimeoutError) Error() string {
	return fmt.Sprintf("DFS: Write access to filename [%s] has timed out; reopen the file", string(e))
}

// Contains filename
type BadFilenameError string

func (e BadFilenameError) Error() string {
	return fmt.Sprintf("DFS: Filename [%s] includes illegal characters or has the wrong length", string(e))
}

// Contains filename
type FileUnavailableError string

func (e FileUnavailableError) Error() string {
	return fmt.Sprintf("DFS: Filename [%s] is unavailable", string(e))
}

// Contains local path
type LocalPathError string

func (e LocalPathError) Error() string {
	return fmt.Sprintf("DFS: Cannot access local path [%s]", string(e))
}

// Contains local path
type NotImplementedError string

func (e NotImplementedError) Error() string {
	return fmt.Sprintf("DFS: The following function has not been implemented: %s\n", string(e))
}

//==================================================================
// The DFS library exposes an interface to the server. The server
// may invoke this interface to request information and data such
// as files stored locally on client or a specific chunk value.
// Client and server communicate via bi-directional RPC calls
//==================================================================

type ClientRPC int

type ClientInterface interface {
	Ping(stub int, reply *bool) (err error)
}
