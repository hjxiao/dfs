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
	myAddr     = "127.0.0.1:3001"
	path       = "../tmp/" // The file path is specified relative to dfslib.go
)

func main() {
	dfs, err := dfslib.MountDFS(serverAddr, myAddr, path)
	exitOnError(err)

	bool1, err := dfs.LocalFileExists("file.dfs")
	fmt.Println("File <file> exists, ", bool1)
	bool2, err := dfs.LocalFileExists("picture.dfs")
	fmt.Println("File <picture> exists, ", bool2)
	bool3, err := dfs.LocalFileExists("gibberish.dfs")
	fmt.Println("File <gibberish> exists, ", bool3)

	f, err := dfs.Open("openTest", 1)
	exitOnError(err)
	var c Chunk
	const str = "Hello world!"
	copy(c[:], str)
	err = f.Write(0, &c)
	fmt.Println(err.Error())

	time.Sleep(10000 * time.Millisecond)

	err = f.Close()
	exitOnError(err)

	time.Sleep(7500 * time.Millisecond)
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
