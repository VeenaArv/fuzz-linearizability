package rqlite

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// For now, all operations must succed or else program crashes.
// TODO(VeenaArv): Consider adding logging.
func runOperation(input string, table *Table, history chan string, wg *sync.WaitGroup) {
	// fmt.Println()
	// fmt.Printf("runOperation %s\n", input)
	inputArr := strings.Split(input, " ")
	// if len(inputArr) < 2 {
	// 	wg.Done()
	// 	return
	// }
	pid := inputArr[0]
	op := inputArr[1]
	// fmt.Println(op)
	if op == "Read" {
		// TODO(veena): MAKE NON-BLOCKING
		history <- fmt.Sprintf("%s Call Read", pid)
		val, err := table.Read()
		if err != nil {
			panic(err)
		}
		history <- fmt.Sprintf("%s Return Read %d", pid, val)

	} else {
		v := strings.Split(inputArr[2], "\r")[0]
		val, err := strconv.ParseInt(v, 10, 0)
		ival := int(val)
		// Success is always true if err is nil.
		// TODO(VeenaArv) Consider removing success from Write.
		if err != nil {
			panic(err)
		}
		// fmt.Printf("%s Call Write %d", pid, val)
		history <- fmt.Sprintf("%s Call Write %d", pid, ival)
		_, err = table.Write(ival)
		if err != nil {
			panic(err)
		}
		history <- fmt.Sprintf("%s Return Write", pid)
		// fmt.Printf("%s Return Write", pid)
	}
	wg.Done()
}

func writeHistory(history chan string, done chan bool, filePath string) {
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	// c := 0
	for true {
		log := <-history
		fmt.Println(log)
		// c++
		if log == "done" {
			err = f.Close()
			if err != nil {
				fmt.Println(err)
				done <- false
				return
			}
			done <- true
			return
		}
		// fmt.Printf("history %s\n", log)
		_, err = fmt.Fprintln(f, log)
	}
	// for d := range data {
	// 	_, err = fmt.Fprintln(f, d)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		f.Close()
	// 		done <- false
	// 		return
	// 	}
	// }
}

func RunOperations(input string, filePath string, strongReadConsistency bool, delays bool) {
	lines := strings.Split(input, "\n")
	numProcesses := 5

	table := NewTable(4001, "test", strongReadConsistency)
	table.CreateTable()

	// history contains call and return logs to feed into linearizability checker.
	history := make(chan string)
	done := make(chan bool)
	channels := make([]chan string, numProcesses)
	wg := sync.WaitGroup{}
	// fmt.Println(numProcesses)
	for i := 0; i < numProcesses; i++ {
		channels[i] = make(chan string)
	}
	for pid, channel := range channels {
		go worker(pid, channel, table, history, &wg, delays)
	}
	go writeHistory(history, done, filePath)
	for i := 1; i < len(lines); i++ {
		wg.Add(1)
		pid, _ := strconv.Atoi(strings.Split(lines[i], " ")[0])
		channels[pid-1] <- lines[i]
	}

	go func() {
		wg.Wait()
		history <- "done"
		close(history)
	}()
	d := <-done
	if d == true {
		// fmt.Println("File written successfully")
	} else {
		fmt.Println("File writing failed")
	}
	table.DeleteTable()
}

func worker(pid int, channel chan string, table *Table, history chan string, wg *sync.WaitGroup, delays bool) {
	// c := 0
	for true {
		// c++
		input := <-channel
		if delays {
			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		}
		// fmt.Printf("input %s\n", input)
		inputArr := strings.Split(input, " ")
		if len(inputArr) < 2 {
			continue
		} else {
			runOperation(input, table, history, wg)
		}
	}
	// fmt.Println(c)
}
