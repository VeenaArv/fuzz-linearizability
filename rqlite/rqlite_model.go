package rqlite

import (
	"github.com/anishathalye/porcupine"
)

type Register struct {
	op    string
	value int
}

// SingleColumnModel is taken from github.com/anishathalye/porcupine/README.md without modification.
func SingleColumnModel() porcupine.Model {
	// a sequential specification of a register
	return porcupine.Model{
		Init: func() interface{} {
			return 0
		},
		// step function: takes a state, input, and output, and returns whether it
		// was a legal operation, along with a new state
		Step: func(state, input, output interface{}) (bool, interface{}) {
			regInput := input.(Register)
			if regInput.op == "read" {
				return true, regInput.value // always ok to execute a write
			}
			readCorrectValue := output == state
			return readCorrectValue, state // state is unchanged
		},
	}

}

func makeSingleColumnCallEvent(op string, value, int, nodeID int) porcupine.Event {
	return porcupine.Event{Kind: porcupine.CallEvent, Value: Register{op, value}, Id: nodeID, ClientId: nodeID}
}

func makeSingleColumnReturnEvent(op string, value int, nodeID int) porcupine.Event {
	return porcupine.Event{Kind: porcupine.ReturnEvent, Value: value, Id: nodeID, ClientId: nodeID}

}

// TODO(VeenaArv): add other models for rqlite.
