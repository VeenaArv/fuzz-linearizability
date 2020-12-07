package rqlite

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/anishathalye/porcupine"
)

type Register struct {
	op    string
	value int
}

// ReadWriteModel is taken from github.com/anishathalye/porcupine/README.md without modification.
func ReadWriteModel() porcupine.Model {
	// a sequential specification of a register
	return porcupine.Model{
		Init: func() interface{} {
			return 0
		},
		// step function: takes a state, input, and output, and returns whether it
		// was a legal operation, along with a new state
		Step: func(state, input, output interface{}) (bool, interface{}) {
			regInput := input.(Register)
			if regInput.op == "Write" {
				return true, regInput.value // always ok to execute a write
			}

			readCorrectValue := output == state
			return readCorrectValue, state // state is unchanged
		},
	}

}

func makeEvent(log string, id int) porcupine.Event {
	tokens := strings.Split(log, " ")
	clientID, _ := strconv.Atoi(tokens[0])
	opType := tokens[1]
	op := tokens[2]
	value := 0
	if len(tokens) == 4 {
		value, _ = strconv.Atoi(tokens[3])
	}
	if opType == "Call" {
		return porcupine.Event{Kind: porcupine.CallEvent, Value: Register{op, value}, Id: id, ClientId: clientID}
	}
	return porcupine.Event{Kind: porcupine.ReturnEvent, Value: value, Id: id, ClientId: clientID}

}

func getEventId(event string, pendingEvents map[string]int, currentID *int) int {
	// truncate value from events if exists.
	tokens := strings.Split(event, " ")[:3]
	if tokens[1] == "Call" {
		*currentID = *currentID + 1
		pendingEvents[strings.Join(tokens, " ")] = *currentID
		return *currentID
	}
	callEvent := tokens[0] + " Call " + tokens[2]
	id := pendingEvents[callEvent]
	delete(pendingEvents, callEvent)
	return id
}

func CheckHistory(filePath string, delFile bool) bool {
	// Reads History
	data, _ := ioutil.ReadFile(filePath)
	lines := strings.Split(string(data), "\n")
	var events []porcupine.Event
	var pendingEvents map[string]int = make(map[string]int)
	currentID := -1
	for i := 0; i < len(lines)-1; i++ {
		line := lines[i]
		// fmt.Println(line)
		id := getEventId(line, pendingEvents, &currentID)
		events = append(events, makeEvent(line, id))
		// fmt.Println(events)
	}
	fmt.Println(events)
	// CheckEvents
	ok := porcupine.CheckEvents(ReadWriteModel(), events)
	if delFile {
		os.Remove(filePath)
	}
	return ok
}

func TestHistory() bool {
	events := []porcupine.Event{
		// 0 call write 1
		{0, porcupine.CallEvent, Register{"Write", 1}, 0},
		// 1 call write 3
		{1, porcupine.CallEvent, Register{"Write", 3}, 1},
		// 0 return write
		{0, porcupine.ReturnEvent, 0, 0},
		// 1 return write
		{1, porcupine.ReturnEvent, 0, 1},
		// 1 call read
		{1, porcupine.CallEvent, Register{"Read", 0}, 2},
		// 1 read 1 or 1 read 3 both return true ; everythin else is false
		{1, porcupine.ReturnEvent, 1, 2},
	}
	// fmt.Println(events)
	return porcupine.CheckEvents(ReadWriteModel(), events)
}

// TODO(VeenaArv): add other models for rqlite.
