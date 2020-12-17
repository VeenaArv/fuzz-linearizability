package fuzzing

import (
	"errors"
	"math/rand"
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

type AlgoRunParams struct {
	NumEvents             int
	NumTests              int
	Run                   int
	StrongReadConsistency bool
	Delays                bool
	Version               string // Used to organized output files to type of algo.
}

func createEvents(numEvents int) []Event {
	rand.Seed(time.Now().Unix())
	events := []Event{}

	ops := []string{"Write", "Read"}

	for i := 0; i < numEvents; i++ {
		op := ops[rand.Intn(2)]
		val := rand.Intn(50)
		pid := rand.Intn(4) + 1 // 1-based
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

func eventsToString(events []Event) string {
	content := numProcesses(events) + "\n"
	for _, e := range events {
		content += strconv.Itoa(e.pid) + " " + e.op
		if e.op == "Write" {
			content += " " + strconv.Itoa(e.val)
		}
		content += "\n"
	}
	return strings.TrimSpace(content)
}

func flipEvent(events []Event, index int) []Event {
	if events[index].op == "Write" {
		events[index].op = "Read"
	} else {
		events[index].op = "Write"
	}
	return events
}

func numProcesses(events []Event) string {
	maxPid := 1
	for _, e := range events {
		if e.pid > maxPid {
			maxPid = e.pid
		}
	}
	return strconv.Itoa(maxPid)
}

func GeneticAlgoWithIncreasingTestCases(params AlgoRunParams) {
	// find optimal value of numEvents to get linearizability errors.
	for numEvents := 10; numEvents < 30; numEvents++ {
		params = AlgoRunParams{numEvents, params.NumTests, params.Run, params.StrongReadConsistency, params.Delays, params.Version}
		if GeneticAlgo(params) {
			break
		}
	}
}

func GeneticAlgo(params AlgoRunParams) bool {
	id := 0
	var events []Event
	var prevEvents = [][]Event{}
	prevIsLinearizable := true
	foundNonLinearizableEvents := false
	flipIndex := 0
	for id < params.NumTests {
		// filePath := fmt.Sprintf("%s/history_%d.txt", historyDirPath, id)
		if prevIsLinearizable {
			events = createEvents(params.NumEvents)
			flipIndex = 0
		} else {
			prevOp := events[flipIndex].op
			events = flipEvent(events, flipIndex)
			if prevOp == events[flipIndex].op {
				panic(errors.New("index not flipped"))
			}
			flipIndex++
			if flipIndex == params.NumEvents {
				prevIsLinearizable = true
			}
		}
		diff := compareEvents(events, prevEvents)
		if !diff {
			continue
		}
		prevEvents = append(prevEvents, events)
		testCaseStats := runTestCase(events, params, id)
		prevIsLinearizable = testCaseStats.linearizable
		if !prevIsLinearizable {
			foundNonLinearizableEvents = true
		}
		id++
	}
	return foundNonLinearizableEvents
}

func RandomizedTesting(params AlgoRunParams) {
	prevEvents := [][]Event{}
	id := 0
	for id < params.NumTests {
		events := createEvents(params.NumEvents)
		diff := compareEvents(events, prevEvents)
		prevEvents = append(prevEvents, events)
		if diff {
			runTestCase(events, params, id)
			id++
		}
	}
}

func runTestCase(events []Event, params AlgoRunParams, id int) TestCaseStats {
	content := eventsToString(events)
	// fmt.Println(content)
	testCaseStats := CheckLinearizability(content, params, id)
	WriteStats(testCaseStats, params, id)
	return testCaseStats
}
