/*
 To run, first install go-fuzz using the commands below:

$ go get github.com/dvyukov/go-fuzz/go-fuzz
$ go get github.com/dvyukov/go-fuzz/go-fuzz-build


Fuzz tells go-fuzz what an interesting input is. Input must contains 1 or
more lines with go-fuzz only considers an input as interesting if there is
new coverage.

Note, for linearizability errors, this approach may not work because
we are concerned with an invalid history. Thus, a input should be considered
interesting if it leads to a new history.

*/

package fuzzers

import (
	"strconv"
	"strings"
)

func isValidInput(s string) int {
	lines := strings.Split(s, "\n")
	numProcesses, err := strconv.Atoi(lines[0])
	if err != nil {
		return -1
	}

	for i := 1; i < len(lines); i++ {
		ret := processSingleLine(lines[i], numProcesses)
		if ret == -1 {
			return -1
		}
	}
	return 1
}
func processSingleLine(s string, numProcesses int) int {
	sArr := strings.Split(s, " ")
	pid, err := strconv.Atoi(sArr[0])
	if err != nil || pid > numProcesses {
		return -1
	}
	// pid Read
	if len(sArr) == 2 && sArr[1] == "Read" {
		return 1
	}
	// pid Write int
	if len(sArr) == 3 && sArr[1] == "Write" {
		_, err := strconv.Atoi(sArr[2])
		if err == nil {
			return 1
		}
	}
	return -1
}
