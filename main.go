package main

import (
	"fmt"
	"fuzz-linearizability/rqlite"
	"io/ioutil"
)

func main() {
	files := []string{"go_fuzz_integration/corpus/input1.txt",
		"go_fuzz_integration/corpus/input2.txt",
		"go_fuzz_integration/corpus/input3.txt"}
	for i, file := range files {
		filePath := fmt.Sprintf("output/histories/history_%d.txt", i)
		content, _ := ioutil.ReadFile(file)
		// This applies operations to rqlite and writes history to filePath.
		rqlite.RunOperations(string(content), filePath, false /*strongReadConsistency*/)
		// This uses porcupine to check the history in filePath and returns
		// true if linearizable.
		fmt.Println(rqlite.CheckHistory(filePath, false /*delFile*/))
	}
	// fmt.Println(rqlite.TestHistory())

}
