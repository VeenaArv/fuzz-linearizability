package rqlite

import (
	"fmt"
	"io/ioutil"
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
	fmt.Println("hi")
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

func makeEvent(log string) porcupine.Event {
	tokens := strings.Split(log, " ")
	nodeID, _ := strconv.Atoi(tokens[0])
	opType := tokens[1]
	op := tokens[2]
	value := 0
	if len(tokens) == 4 {
		value, _ = strconv.Atoi(tokens[3])
	}
	if opType == "Call" {
		return porcupine.Event{Kind: porcupine.CallEvent, Value: Register{op, value}, Id: nodeID, ClientId: nodeID}
	}
	return porcupine.Event{Kind: porcupine.ReturnEvent, Value: value, Id: nodeID, ClientId: nodeID}

}

func CheckHistory() bool {
	// Reads History
	data, _ := ioutil.ReadFile("output/history.txt")
	lines := strings.Split(string(data), "\n")
	var events []porcupine.Event
	for i := 0; i < len(lines)-1; i++ {
		line := lines[i]
		fmt.Println(line)
		events = append(events, makeEvent(line))
		// fmt.Println(events)
	}
	fmt.Println(events)
	// CheckEvents
	ok := porcupine.CheckEvents(ReadWriteModel(), events)
	// fmt.Println(res)
	// fmt.Println(info)
	return ok
	// data := computeVisualizationData(ReadWriteModel(), info)
}

// TODO(VeenaArv): add other models for rqlite.
