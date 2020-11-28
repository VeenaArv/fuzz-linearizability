package main

import (
	"fmt"
	// "io/ioutil"
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
	endpoint string
	name     string
}

// NewTable creates and returns a pointer to an instance of an Table.
func NewTable(port int, name string) *Table {
	endpoint := fmt.Sprintf("localhost:%d/db/query?level=strong", port)
	return &Table{endpoint, name}
}

func (table Table) Read() {
	// Always enable strong consistency which garuntees linearizability. See
	// https://github.com/rqlite/rqlite/blob/master/DOC/CONSISTENCY.md for
	//different consistency garuntees.
	query := fmt.Sprintf("SELECT * FROM %s", table.name)
	runQuery(table, query, "GET")
}

// TODO(veenaarv): update this function to insert values of any type and
// support multiple columns.
// Table must be created using `CreateTable`, before writing to it.
func (table Table) Write(value int) {
	query := fmt.Sprintf("INSERT INTO %s(value) VALUES(%d)", table.name, value)
	runQuery(table, query, "POST")
}

func (table Table) CreateTable() {
	query := fmt.Sprintf("CREATE TABLE %s (value INTEGER)", table.name)
	runQuery(table, query, "POST")
}

func runQuery(table Table, query string, httpMethod string) *http.Request {
	req, _ := http.NewRequest(httpMethod, table.endpoint, nil)
	q := req.URL.Query()          // Get a copy of the query values.
	q.Add("q", query)             // Add a new value to the set.
	req.URL.RawQuery = q.Encode() // Encode and assign back to the original query.
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	fmt.Printf("URL      %+v\n", req.URL)
	fmt.Printf("RawQuery %+v\n", req.URL.RawQuery)
	fmt.Printf("Query    %+v\n", req.URL.Query())

	return req
}
func main() {
	fmt.Println("Hello, world.")

}
