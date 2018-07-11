package main

import (
	"fmt"
	"os"
)

var (
	fd *os.File
)

func main() {
	fd, err := os.Create("file.dfs")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println(fd.Name())
	buf := make([]byte, 32)
	fd.Truncate(int64(len(buf) * 256))
	fd.Seek(0, 0)
	count, err := fd.Read(buf)
	fmt.Printf("Read %d bytes\n", count)
	fmt.Println(buf)

	buf2 := make([]byte, 32)
	const str = "Hello world!"
	copy(buf2, str)
	fmt.Println("buf2: ", buf2)

	fd.Seek(0, 0)
	count2, err := fd.Write(buf2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Wrote %d bytes\n", count2)
	err = fd.Sync()
	if err != nil {
		fmt.Println(err)
	}

	buf3 := make([]byte, 32)
	fd.Seek(0, 0)
	count, err = fd.Read(buf3)
	fmt.Printf("Read %d bytes\n", count)
	fmt.Println(buf3)

	fd.Close()
}
