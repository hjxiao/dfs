package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	files       = make(map[string]FileState, 0)
	writeAccess = make(map[string]sync.Mutex, 0)
)

type FileState struct {
	fileExists   bool
	writeLock    *sync.Mutex
	chunkVersion []*FileVersionOwners // All chunks initialized at version 0; each write increments by 1
}

type FileVersionOwners struct {
	version int
	owners  []UserInfo
}

type UserInfo struct {
	LocalIP   string
	LocalPath string
}

func main() {
	fs1 := FileState{fileExists: true, writeLock: &sync.Mutex{}, chunkVersion: make([]*FileVersionOwners, 256)}
	files["test1"] = fs1

	go func() {
		fmt.Println("fnc1: Attempting to acquire test1 lock")
		files["test1"].writeLock.Lock()
		fmt.Println("fnc1: Acquired lock")
		fmt.Println("fnc1: Sleeping for 5 seconds")
		time.Sleep(5 * time.Second)
		fmt.Println("fnc1: Waking and releasing lock")
		files["test1"].writeLock.Unlock()
		fmt.Println("fnc1: Lock released")
	}()

	go func() {
		fmt.Println("fnc2: Attempting to acquire test1 lock")
		files["test1"].writeLock.Lock()
		fmt.Println("fnc2: Acquired lock")
		fmt.Println("fnc2: Sleeping for 5 seconds")
		time.Sleep(5 * time.Second)
		fmt.Println("fnc2: Waking and releasing lock")
		files["test1"].writeLock.Unlock()
		fmt.Println("fnc2: Lock released")
	}()

	time.Sleep(10000 * time.Second)
}
