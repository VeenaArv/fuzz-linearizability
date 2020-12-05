package rqlite

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type TestProcess struct {
	pid int
	op  string // Must be either "Read", "Write X"
}

// For now, all operations must succed or else program crashes.
// TODO(VeenaArv): Consider adding logging.
func runOperation(process TestProcess, table Table, opLog chan string, wg *sync.WaitGroup) {
	if process.op == "Read" {
		opLog <- fmt.Sprintf("%d Call Read", process.pid)
		val, err := table.Read()
		if err != nil {
			panic(err)
		}
		opLog <- fmt.Sprintf("%d Return Read %d", process.pid, val)

	} else {
		val, err := strconv.Atoi(strings.Split(process.op, " ")[1])
		// Success is always true if err is nil.
		// TODO(VeenaArv) Consider removing success from Write.
		if err != nil {
			panic(err)
		}
		opLog <- fmt.Sprintf("%d Call Write %d", process.pid, val)
		_, err = table.Write(val)
		if err != nil {
			panic(err)
		}
		opLog <- fmt.Sprintf("%d Return Write", process.pid)
	}
	wg.Done()
}

func writeHistory(data chan int, done chan bool) {
	f, err := os.Create("history_test_input.txt")
	if err != nil {
		panic(err)
	}
	for d := range data {
		_, err = fmt.Fprintln(f, d)
		if err != nil {
			fmt.Println(err)
			f.Close()
			done <- false
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		done <- false
		return
	}
	done <- true
}

func runOperations(input string) {
	data := make(chan string)
	done := make(chan bool)
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		// go runOperation
	}

	// go writeHistory
	go func() {
		wg.Wait()
		close(data)
	}()
	d := <-done
	if d == true {
		fmt.Println("File written successfully")
	} else {
		fmt.Println("File writing failed")
	}
}
