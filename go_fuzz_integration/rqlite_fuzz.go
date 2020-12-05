package fuzzers

import (
	"fuzz-linearizability/rqlite"
)

/*
 To run, first install go-fuzz using the commands below:

$ go get github.com/dvyukov/go-fuzz/go-fuzz
$ go get github.com/dvyukov/go-fuzz/go-fuzz-build


Fuzz tells go-fuzz what an interesting input is. Input must contains 1 or
more lines with go-fuzz only consideres an input as interesting if there is
new coverage. For linearizability errors, this approach may not work because
we are concerned with an invalid history. Thus, a input should be considered
interesting if it leads to a new history.

*/

// Fuzz input must be of form `pid Read` or `pid Write val`
// Non-Linearizable histories are considered most important
func Fuzz(data []byte) int {
	// if _, err := Decode(bytes.NewReader(data)); err != nil {
	s := string(data)
	isValid := isValidInput(s)
	t := rqlite.NewTable(4001, "test")

}
