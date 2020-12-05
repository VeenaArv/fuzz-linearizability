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
