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

	bool1, err := dfs.LocalFileExists("file")
	fmt.Println("File <file> exists, ", bool1)
	bool2, err := dfs.LocalFileExists("picture")
	fmt.Println("File <picture> exists, ", bool2)
	bool3, err := dfs.LocalFileExists("gibberish")
	fmt.Println("File <gibberish> exists, ", bool3)

	f, err := dfs.Open("openTest", 2)
	exitOnError(err)

	var c Chunk
	const str = "Hello world!"
	copy(c[:], str)
	fmt.Println("Writing chunk to file: ", c)
	err = f.Write(3, &c)
	exitOnError(err)

	var c2 Chunk
	err = f.Read(3, &c2)
	exitOnError(err)
	fmt.Println("Read chunk from file: ", c2)

	var c3 Chunk
	const str2 = "Testing stale read"
	copy(c3[:], str2)
	fmt.Println("Writing chunk to file: ", c3)
	err = f.Write(10, &c3)
	err = f.Write(10, &c3)
	exitOnError(err)

	time.Sleep(20000 * time.Millisecond)

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
