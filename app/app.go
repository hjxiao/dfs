package main

import (
	"fmt"
	"os"

	"../dfslib"
)

var (
	serverAddr = "127.0.0.1:3000"
	myAddr     = "127.0.0.1:3001"
	path       = "/tmp/"
)

func main() {
	dfs, err := dfslib.MountDFS(serverAddr, myAddr, path)
	exitOnError(err)

	err = dfs.UMountDFS()

	for {

	}
}

func exitOnError(e error) {
	if e != nil {
		fmt.Println("app: Encountered error - ", e.Error())
		os.Exit(-1)
	}
}
