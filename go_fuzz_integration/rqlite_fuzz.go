package fuzzers

import (
	"fuzz-linearizability/rqlite"
)

// Fuzz input must be of form `pid Read` or `pid Write val`
// Non-Linearizable histories are considered most important
func Fuzz(data []byte) int {
	// if _, err := Decode(bytes.NewReader(data)); err != nil {
	input := string(data)
	isValid := isValidInput(input)
	t := rqlite.NewTable(4001, "test")
	rqlite.RunOperations(input)
	linearizable := rqlite.CheckHistory("output/history.txt")
	if linearizable {
		return 1
	}
	return 0
	// TODO(veena)

}
