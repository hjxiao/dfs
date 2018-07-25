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
	myAddr     = "127.0.0.1:3005"
	path       = "../tmp2/" // The file path is specified relative to dfslib.go
)

func main() {
	dfs, err := dfslib.MountDFS(serverAddr, myAddr, path)
	exitOnError(err)

	f, err := dfs.Open("secondDRTest", 2)
	exitOnError(err)

	var c Chunk
	const str = "Fresh data for everyone!"
	copy(c[:], str)
	fmt.Println("Writing chunk to file: ", c)
	err = f.Write(7, &c)
	exitOnError(err)

	time.Sleep(20000 * time.Millisecond)

	err = f.Close()
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
