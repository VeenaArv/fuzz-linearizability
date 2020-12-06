package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Event is a struct
type Event struct {
	op  string
	val int
	pid int
}

// create events
func createEvents() []Event {
	events := []Event{}

	ops := []string{"Write", "Read"}

	for i := 0; i < 5; i++ {
		op := ops[rand.Intn(2)]
		val := rand.Intn(50)
		pid := rand.Intn(5)
		events = append(events, Event{op, val, pid})
	}

	return events
}

// write function
func write(toWrite int, pid int) {
	// fmt.Println("HEREEE")
	object := strconv.Itoa(pid) + " " + strconv.FormatInt(time.Now().Unix(), 10)
	// err := ioutil.WriteFile(object, []byte(strconv.Itoa(toWrite)), 0644)
	f, err := os.Create(object)
	fmt.Println(err)
	f.WriteString(strconv.Itoa(toWrite) + "\n")
	upload(toWrite, object)
	// defer f.Close()
}

// read function
func read() {
	fmt.Println("read")
	// download()
}

func worker(channel chan Event, pid int) {
	for {
		command := <-channel

		if command.op == "Write" {
			write(command.val, pid)
		} else if command.op == "Read" {
			read()
		}
	}
}

// feed events
func applyEvents(events []Event) {
	// add fuzzer
	// get history
	pids := []int{0, 1, 2, 3, 4}
	channels := []chan Event{}

	for i := range pids {
		newChannel := make(chan Event)
		channels = append(channels, newChannel)
		go worker(newChannel, i)
	}

	for _, event := range events {
		channels[event.pid] <- event
	}
}

// feed history to porcupine
func testHistory(history string) {
	// TODO: evaluate
	// TODO: GENETIC

	// porcupine.CheckEvents
}

func getHistory() []byte {
	data, _ := ioutil.ReadFile("../data/sample_history.txt")
	// fmt.Println(data)
	return data
}

// check order of reads and writes in history for comparison
func parseHistory(h string) string {
	order := ""
	words := strings.Split(h, " ")

	for i := 2; i < len(words); i += 4 {
		order = order + words[i]
	}
	return order
}

func compareHistory(history string, prevHistories []string) bool {
	for _, h := range prevHistories {
		// fmt.Println(h)
		if parseHistory(h) == parseHistory(history) {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println("Hello World")

	rand.Seed(time.Now().Unix())

	prevHistories := []string{}
	for i := 0; i < 10; i++ {
		events := createEvents()
		applyEvents(events)
		history := string(getHistory())
		prevHistories = append(prevHistories, history)
		diff := compareHistory(history, prevHistories)
		if diff {
			testHistory(history)
		}
		// else {
		// 	i--
		// }
		fmt.Println(diff)

	}
	// testHistory(history)

	// fmt.Println(events)
}
