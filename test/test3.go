package main

import (
	"os"
	"unicode"
)

func validFileName(str string) bool {
	if len(str) < 1 || len(str) > 16 {
		return false
	}

	return isAlphaNumeric(str)
}

func isAlphaNumeric(str string) bool {
	for _, s := range str {
		if !unicode.IsLetter(s) && !unicode.IsNumber(s) {
			return false
		}
	}
	return true
}

func main() {
	// fmt.Println(validFileName("hello123"))
	// fmt.Println(validFileName(""))
	// fmt.Println(validFileName("&^"))
	// fmt.Println(validFileName("hello"))

	fd, err := os.Create("../tmp/random.dfs")
	if err != nil {
		panic(err)
	}

	fd.Close()
}
