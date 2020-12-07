package main

import (
	"fmt"
	"fuzz-linearizability/rqlite"
	"math/rand"
	"strconv"
)

// Event is a struct
type Event struct {
	op  string
	val int
	pid int
}

func createEvents() []Event {
	events := []Event{}

	ops := []string{"Write", "Read"}

	for i := 0; i < 8; i++ {
		op := ops[rand.Intn(2)]
		val := rand.Intn(50)
		pid := rand.Intn(5)
		events = append(events, Event{op, val, pid})
	}

	return events
}

// returns true if two sets of events are different
func diffEvents(e1 []Event, e2 []Event) bool {
	for i := 0; i < len(e1); i++ {
		if e1[i].op != e2[i].op {
			return true
		}
	}
	return false
}

func compareEvents(events []Event, prevEvents [][]Event) bool {
	for _, e := range prevEvents {
		if diffEvents(events, e) == false {
			return false
		}
	}
	return true
}

func numProcesses(events []Event) string {
	active := []int{0, 0, 0, 0, 0}
	for _, e := range events {
		active[e.pid] = 1
	}
	c := 0
	for _, p := range active {
		if p == 1 {
			c++
		}
	}
	return strconv.Itoa(c)
}

func main() {
	files := []string{"go_fuzz_integration/corpus/input1.txt",
		"go_fuzz_integration/corpus/input2.txt",
		"go_fuzz_integration/corpus/input3.txt"}

	// non-repeating events, network delays, mutation instead of random
	// non-repeating histories rquires running the events so test cases arent decreased no?
	strategies := []int{0, 1, 2}

	for s := range strategies {
		switch s {
		case 0:
			prevEvents := [][]Event{}
			for i := 0; i < 10; i++ {
				events := createEvents()
				prevEvents = append(prevEvents, events)
				diff := compareEvents(events, prevEvents)
				if diff {
					filePath := fmt.Sprintf("output/histories/history_%d.txt", i)
					content := ""
					content = content + numProcesses(events) + "\n"
					for _, e := range events {
						content = content + strconv.Itoa(e.pid) + " " + e.op + " " + strconv.Itoa(e.val) + "\n"
					}
					rqlite.RunOperations(string(content), filePath, false /*strongReadConsistency*/, false /*delays*/)
					fmt.Println(rqlite.CheckHistory(filePath, false /*delFile*/))
				}
				// else {
				// 	i--
				// }
			}
		// case 1:
		// 	// prevHistories := []string{}
		// 	// for i := 0; i < 10; i++ {
		// 	// 	events := createEvents()
		// 	// 	filePath := fmt.Sprintf("output/histories/history_%d.txt", i)
		// 	// 	content := ""
		// 	// 	content = content + numProcesses(events) + "\n"
		// 	// 	for _, e := range events {
		// 	// 		content = content + strconv.Itoa(e.pid) + " " + e.op + " " + strconv.Itoa(e.val) + "\n"
		// 	// 	}
		// 	// 	rqlite.RunOperations(string(content), filePath, false /*strongReadConsistency*/, false /*delays*/)
		// 	// 	fmt.Println(rqlite.CheckHistory(filePath, false /*delFile*/))
		// 	// 	diff := compareHistory(history, prevHistories)
		// 	// 	if diff {
		// 	// 		testHistory(history)
		// 	// 	}
		// 		// else {
		// 		// 	i--
		// 		// }
		// 		// fmt.Println(diff)
		// 	// }
		case 1:
			for i := 0; i < 10; i++ {
				events := createEvents()
				filePath := fmt.Sprintf("output/histories/history_%d.txt", i)
				content := ""
				content = content + numProcesses(events) + "\n"
				for _, e := range events {
					content = content + strconv.Itoa(e.pid) + " " + e.op + " " + strconv.Itoa(e.val) + "\n"
				}
				rqlite.RunOperations(string(content), filePath, false /*strongReadConsistency*/, true /*delays*/)
				fmt.Println(rqlite.CheckHistory(filePath, false /*delFile*/))
			}
		case 2:
			for i := 0; i < 10; i++ {
				events := createEvents()
				filePath := fmt.Sprintf("output/histories/history_%d.txt", i)
				content := ""
				content = content + numProcesses(events) + "\n"
				for _, e := range events {
					content = content + strconv.Itoa(e.pid) + " " + e.op + " " + strconv.Itoa(e.val) + "\n"
				}
				rqlite.RunOperations(string(content), filePath, false /*strongReadConsistency*/, false /*delays*/)
				h := rqlite.CheckHistory(filePath, false /*delFile*/)
				fmt.Println(h)
				if h == false {
					for j := 0; j < 8; j++ {
						new_events := events
						if new_events[j].op == "Write" {
							new_events[j].op == "Read"
						} else {
							new_events[j].op == "Write"
						}
						content := ""
						content = content + numProcesses(new_events) + "\n"
						for _, e := range new_events {
							content = content + strconv.Itoa(e.pid) + " " + e.op + " " + strconv.Itoa(e.val) + "\n"
						}
						rqlite.RunOperations(string(content), filePath, false /*strongReadConsistency*/, false /*delays*/)
						fmt.Println(rqlite.CheckHistory(filePath, false /*delFile*/))
					}
					break
				}
			}
		default:
			continue
		}
	}

	// for i, file := range files {
	// filePath := fmt.Sprintf("output/histories/history_%d.txt", i)
	// content, _ := ioutil.ReadFile(file)
	// This applies operations to rqlite and writes history to filePath.
	// rqlite.RunOperations(string(content), filePath, false /*strongReadConsistency*/)
	// This uses porcupine to check the history in filePath and returns
	// true if linearizable.
	// fmt.Println(rqlite.CheckHistory(filePath, false /*delFile*/))
	// }
	// fmt.Println(rqlite.TestHistory())

}
