package main

import (
	"fmt"
	"math/rand"

	// "github.com/anishathalye/porcupine";
	"context"
	"io"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

/*
TODO:
1. figure out genetic algo fitness function
2. feed generated trace to gcp
3. feed history to porcupine
*/

// uploadFile uploads an object.
func uploadFile(w io.Writer, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open("notes.txt")
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Fprintf(w, "Blob %v uploaded.\n", object)
	return nil
}

// downloadFile downloads an object.
func downloadFile(w io.Writer, bucket, object string) ([]byte, error) {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	fmt.Fprintf(w, "Blob %v downloaded.\n", object)
	return data, nil
}

type registerInput struct {
	op    bool // false = write, true = read
	value int
}

func worker(pid int, channel chan int) {
	for true {
		// event := <-channel
		// if event.op == "w" {
		// 	err := uploadFile(event.val)
		// } else if event.op == "r" {
		// 	val, err := event.val
		// }

		// fmt.Println(<-channel)
		fmt.Println(pid)
		<-channel

	}
}

// Event Type of operation
type Event struct {
	pid       int
	operation string
}

func bestScore(operations []Event) int {
	finalList := []Event{}
	for _, e := range operations {
		score := 0
		if e.operation == "Write" {
			return 2
		} else {
			return 1
		}
	}
}
func generateEvents() []string {
	events := []string{}
	pidOptions := [...]int{0, 1, 2, 3, 4}
	operations := []string{"Write", "Read"}

	possibilities := []Event{}
	for i := 1; i < 20; i++ {
		possibilities = append(possibilities, Event{rand.Intn(5), operations[rand.Intn(2)]})
	}
	bestScore(possibilities)

	return events
}

func main() {
	// trace1 := [4]string{"r", "w", "r", "r"}

	chans := []chan int{
		make(chan int),
		make(chan int),
		make(chan int),
		make(chan int),
		make(chan int),
	}

	for pid, channel := range chans {
		go worker(pid, channel)
	}
	chans[0] <- 1
	chans[1] <- 5
	chans[3] <- 5
	chans[4] <- 2
	// for _, task := range trace1 {
	// 	if task == "r" {
	// 		go downloadFile(os.Stdout, "test_ds_object", "note")
	// 	} else if task == "w" {
	// 		go uploadFile(os.Stdout, "test_ds_object", "note")
	// 	}
	// }

	time.Sleep(3 * time.Second)

	// uploadFile(os.Stdout, "test_ds_object", "note")
	// a sequential specification of a register
	// registerMTrn 0
	// 	},
	// 	// step function: takes a state, input, and output, and returns whether it
	// 	// was a legal operation, along with a new state
	// 	Step: func(state, input, output interface{}) (bool, interface{}) {
	// 		regInput := input.(registerInput)
	// 		if regInput.op == false {
	// 			return true, regInput.value // always ok to execute a write
	// 		} else {
	// 			readCorrectValue := output == state
	// 			return readCorrectValue, state // state is unchanged
	// 		}
	// 	},
	// }

	// events := []porcupine.Event{
	// 	// C0: Write(100)
	// 	{Kind: porcupine.CallEvent, Value: registerInput{false, 100}, Id: 0, ClientId: 0},
	// 	// C1: Read()
	// 	{Kind: porcupine.CallEvent, Value: registerInput{true, 0}, Id: 1, ClientId: 1},
	// 	// C2: Read()
	// 	{Kind: porcupine.CallEvent, Value: registerInput{true, 0}, Id: 2, ClientId: 2},
	// 	// C2: Completed Read -> 0
	// 	{Kind: porcupine.ReturnEvent, Value: 0, Id: 2, ClientId: 2},
	// 	// C1: Completed Read -> 100
	// 	{Kind: porcupine.ReturnEvent, Value: 100, Id: 1, ClientId: 1},
	// 	// Fake read: Completed Read -> 100
	// 	{Kind: porcupine.ReturnEvent, Value: 100, Id: 1, ClientId: 1},
	// 	// C0: Completed Write
	// 	{Kind: porcupine.ReturnEvent, Value: 0, Id: 0, ClientId: 0},
	// }

	// ok := porcupine.CheckEvents(registerModel, events)
	// fmt.Println(ok)
}
