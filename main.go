package main

import (
	"fmt"
	"fuzz-linearizability/rqlite"
	"io/ioutil"
)

func main() {
	content, _ := ioutil.ReadFile("data/sample_input.txt")
	fmt.Println(string(content))
	rqlite.RunOperations(string(content))
	fmt.Println(rqlite.CheckHistory())
}
