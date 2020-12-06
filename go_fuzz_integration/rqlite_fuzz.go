package fuzzers

import (
	"fuzz-linearizability/rqlite"
)

// Fuzz input must be of form `pid Read` or `pid Write val`
// Non-Linearizable histories are considered most important
func Fuzz(data []byte) int {
	// if _, err := Decode(bytes.NewReader(data)); err != nil {
	s := string(data)
	isValid := isValidInput(s)
	t := rqlite.NewTable(4001, "test")
	// TODO(veena)

}
