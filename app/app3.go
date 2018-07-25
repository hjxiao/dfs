package main

import (
	"fmt"
	"os"

	"../dfslib"
)

type Chunk = dfslib.Chunk

var (
	serverAddr = "127.0.0.1:3000"
	myAddr     = "127.0.0.1:3003"
	path       = "../tmp2/" // The file path is specified relative to dfslib.go
)

func main() {
	dfs, err := dfslib.MountDFS(serverAddr, myAddr, path)
	exitOnError(err)

	_, err = dfs.Open("noSuchFile", 3)
	expectingError(err)

	f, err := dfs.Open("dreadTest", 2)
	exitOnError(err)

	var c Chunk
	const str = "Testing disconnected reads"
	copy(c[:], str)
	fmt.Println("Writing chunk to file: ", c)
	err = f.Write(10, &c)
	exitOnError(err)

	err = f.Close()
	exitOnError(err)

	f2, err := dfs.Open("dreadTest", 3)
	exitOnError(err)

	var c2 Chunk
	err = f2.Read(10, &c2)
	exitOnError(err)
	fmt.Println("Disconnected read on chunk from file: ", c2)

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
