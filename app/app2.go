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
	myAddr     = "127.0.0.1:3002"
	path       = "../tmp2/" // The file path is specified relative to dfslib.go
)

func main() {
	dfs, err := dfslib.MountDFS(serverAddr, myAddr, path)
	exitOnError(err)

	bool1, err := dfs.GlobalFileExists("gibberish")
	fmt.Println("File <gibberish> exists, ", bool1)

	bool2, err := dfs.GlobalFileExists("openTest")
	fmt.Println("File <openTest> exists, ", bool2)

	time.Sleep(2000 * time.Millisecond)

	_, err = dfs.Open("openTest", 2)
	fmt.Println("Expecting error: ", err.Error())

	f, err := dfs.Open("openTest", 1)
	exitOnError(err)

	var c Chunk
	err = f.Read(10, &c)
	exitOnError(err)
	fmt.Println("Read chunk from file: ", c)

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
