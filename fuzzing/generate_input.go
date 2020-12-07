package fuzzing

import (
	"io/ioutil"
	"math/rand"
	"strings"
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

// func main() {
// 	fmt.Println("Hello World")

// 	rand.Seed(time.Now().Unix())

// 	prevHistories := []string{}
// 	for i := 0; i < 10; i++ {
// 		events := createEvents()
// 		applyEvents(events)
// 		history := string(getHistory())
// 		prevHistories = append(prevHistories, history)
// 		diff := compareHistory(history, prevHistories)
// 		if diff {
// 			testHistory(history)
// 		}
// 		// else {
// 		// 	i--
// 		// }
// 		fmt.Println(diff)

// 	}
// testHistory(history)

// fmt.Println(events)
// }
