package rqlite

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// For now, all operations must succed or else program crashes.
// TODO(VeenaArv): Consider adding logging.
func runOperation(input string, table *Table, history chan string, wg *sync.WaitGroup) {
	fmt.Printf("runOperation %s\n", input)
	inputArr := strings.Split(input, " ")
	pid := inputArr[0]
	op := inputArr[1]
	fmt.Println(op)
	if op == "Read" {
		fmt.Println("reading")
		// TODO(veena): MAKE NON-BLOCKING
		history <- fmt.Sprintf("%s Call Read", pid)
		val, err := table.Read()
		fmt.Println(val)
		if err != nil {
			panic(err)
		}
		history <- fmt.Sprintf("%s Return Read %d", pid, val)

	} else {
		val, err := strconv.Atoi(inputArr[2])
		fmt.Printf("writing %d", val)
		// Success is always true if err is nil.
		// TODO(VeenaArv) Consider removing success from Write.
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s Call Write %d", pid, val)
		history <- fmt.Sprintf("%s Call Write %d", pid, val)
		_, err = table.Write(val)
		if err != nil {
			panic(err)
		}
		history <- fmt.Sprintf("%s Return Write", pid)
		fmt.Printf("%s Return Write", pid)
	}
	wg.Done()
}

func writeHistory(history chan string, done chan bool) {
	f, err := os.Create("output/history.txt")
	if err != nil {
		panic(err)
	}
	for true {
		log := <-history
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
		fmt.Printf("history %s\n", log)
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

func RunOperations(input string) {
	lines := strings.Split(input, "\n")
	numProcesses, _ := strconv.Atoi(lines[0])

	table := NewTable(4001, "test")
	table.CreateTable()

	// history contains call and return logs to feed into linearizability checker.
	history := make(chan string)
	done := make(chan bool)
	channels := make([]chan string, numProcesses)
	wg := sync.WaitGroup{}
	fmt.Println(numProcesses)
	for i := 0; i < numProcesses; i++ {
		channels[i] = make(chan string)
	}
	for pid, channel := range channels {
		go worker(pid, channel, table, history, &wg)
	}
	fmt.Println("next")
	go writeHistory(history, done)
	for i := 1; i < len(lines); i++ {
		wg.Add(1)
		pid, _ := strconv.Atoi(strings.Split(lines[i], " ")[0])
		fmt.Println(lines[i])
		channels[pid-1] <- lines[i]
	}
	fmt.Println("writing history")

	go func() {
		wg.Wait()
		history <- "done"
		close(history)
	}()
	d := <-done
	if d == true {
		fmt.Println("File written successfully")
	} else {
		fmt.Println("File writing failed")
	}
}

func worker(pid int, channel chan string, table *Table, history chan string, wg *sync.WaitGroup) {
	for true {
		input := <-channel
		fmt.Printf("input %s\n", input)
		runOperation(input, table, history, wg)
	}
}

// func main() {
// 	content, err := ioutil.ReadFile("data.sample_history.txt")
// 	fmt.Println(string(content))
// }
