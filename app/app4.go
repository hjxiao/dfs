package main

import (
	"fmt"
	"os"
	"time"

	"../dfslib"
)

type Chunk = dfslib.Chunk

var (
	serverAddr = "127.0.0.1:3000"
	myAddr     = "127.0.0.1:3004"
	path       = "../tmp/" // The file path is specified relative to dfslib.go
)

func main() {
	dfs, err := dfslib.MountDFS(serverAddr, myAddr, path)
	exitOnError(err)

	f, err := dfs.Open("secondDRTest", 2)
	exitOnError(err)

	var c Chunk
	const str = "Testing stale and fresh read"
	copy(c[:], str)
	fmt.Println("Writing chunk to file: ", c)
	err = f.Write(7, &c)
	exitOnError(err)

	err = f.Close()
	exitOnError(err)

	f2, err := dfs.Open("secondDRTest", 3)
	exitOnError(err)

	var c2 Chunk
	err = f2.Dread(7, &c2)
	exitOnError(err)
	fmt.Println("Stale read on chunk from file: ", c2)

	time.Sleep(10000 * time.Millisecond)

	var c3 Chunk
	err = f2.Read(7, &c3)
	exitOnError(err)
	fmt.Println("Up-to-date read on chunk from file: ", c3)

	err = f2.Close()
	exitOnError(err)

	err = dfs.UMountDFS()
	exitOnError(err)

	for {

	}
}

func exitOnError(e error) {
	if e != nil {
		fmt.Println("app: Encountered error - ", e.Error())
		os.Exit(-1)
	}
}

func expectingError(e error) {
	if e == nil {
		fmt.Println("app: Expected error, but it was nil")
		os.Exit(-1)
	} else {
		fmt.Println("app: Received expected error: ", e.Error())
	}
}
