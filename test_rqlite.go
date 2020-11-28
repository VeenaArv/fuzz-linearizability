package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// RQLiteTest defines operations to test for linearizability.
type RQLiteTest interface {
	Read()
	Write(value int)
	CreateTable()
}

// Table is an RQLite table with name `name` and can be accessed using HTTP rquests at `endpoint`.
// RQLite nodes must be setup before accessing Table.
type Table struct {
	port int
	name string
}

// NewTable creates and returns a pointer to an instance of an Table.
func NewTable(port int, name string) *Table {
	return &Table{port, name}
}

func (table Table) Read() {
	// Always enable strong consistency which garuntees linearizability. See
	// https://github.com/rqlite/rqlite/blob/master/DOC/CONSISTENCY.md for
	//different consistency garuntees.
	// TODO(VeenaArv): add flag to switch between different consistency levels.
	query := fmt.Sprintf("SELECT * FROM %s", table.name)
	getQuery(query, table.port)
}

// TODO(veenaarv): update this function to insert values of any type and
// support multiple columns.
// Table must be created using `CreateTable`, before writing to it.
func (table Table) Write(value int) {
	query := fmt.Sprintf("[\"INSERT INTO %s(value) VALUES(%d)\"]", table.name, value)
	postQuery(query, table.port)
}

func (table Table) CreateTable() {
	query := fmt.Sprintf("[\"CREATE TABLE %s (value INTEGER)\"]", table.name)
	postQuery(query, table.port)
}

func postQuery(query string, port int) *http.Response {
	endpoint := fmt.Sprintf("http://localhost:%d/db/execute?pretty&timings", port)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(query)))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Status)
	return resp
}

func getQuery(query string, port int) *http.Response {
	endpoint := fmt.Sprintf("http://localhost:%d/db/query?pretty&timings", port)
	req, err := http.NewRequest("GET", endpoint, nil)
	q := req.URL.Query() // Get a copy of the query values.

	q.Add("level", "strong")      // read consisteny always strong.
	q.Add("q", query)             // Add a new value to the set.
	req.URL.RawQuery = q.Encode() // Encode and assign back to the original query.
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Status)
	return resp
}
func main() {
	t := NewTable(4001, "test2")
	t.CreateTable()
	// t.Write(1)
	// t.Write(2)
	t.Read()

}
