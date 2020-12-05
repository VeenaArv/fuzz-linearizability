package rqlite

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// RQLiteTest defines operations to test for linearizability.
type RQLiteTest interface {
	Read()
	Write(value int)
	CreateTable()
	DeleteTable()
}

// Table is an RQLite table with name `name` and can be accessed using HTTP rquests at `endpoint`.
// RQLite nodes must be setup before accessing Table.
type Table struct {
	port int
	name string
}

type RQLiteReadJsonResponse struct {
	Results []struct {
		// Columns []string `json:"columns"`
		// Types   []string `json:"types"`
		Values [][]int `json:"values"`
		// Time    float64  `json:"time"`
	} `json:"results"`
	// Time float64 `json:"time"`
}

// NewTable creates and returns a pointer to an instance of an Table.
func NewTable(port int, name string) *Table {
	return &Table{port, name}
}

func (table Table) Read() (int, error) {
	// Always enable strong consistency which garuntees linearizability. See
	// https://github.com/rqlite/rqlite/blob/master/DOC/CONSISTENCY.md for
	//different consistency garuntees.
	// TODO(VeenaArv): add flag to switch between different consistency levels.
	query := fmt.Sprintf("SELECT * FROM %s", table.name)
	endpoint := fmt.Sprintf("http://localhost:%d/db/query?pretty&timings", table.port)
	req, err := http.NewRequest("GET", endpoint, nil)
	q := req.URL.Query()          // Get a copy of the query values.
	q.Add("level", "strong")      // read consisteny always strong.
	q.Add("q", query)             // Add a new value to the set.
	req.URL.RawQuery = q.Encode() // Encode and assign back to the original query.
	resp, err := runQuery(req)
	if err != nil {
		return -1, err
	}
	// Parse response.
	body, err := ioutil.ReadAll(resp.Body)
	json_resp := new(RQLiteReadJsonResponse)
	err = json.Unmarshal(body, &json_resp)
	values := json_resp.Results[0].Values
	if len(values) == 1 && len(values[0]) == 1 {
		return values[0][0], nil
	} else {
		return -1, errors.New("Read does not contain only 1 value.")
	}
}

// TODO(veenaarv): update this function to insert values of any type and
// support multiple columns.
// Table must be created using `CreateTable`, before writing to it.
func (table Table) Write(value int) (bool, error) {

	// updates single row with `value`
	query := fmt.Sprintf("[\"UPDATE %s SET value = %d\"]", table.name, value)
	endpoint := fmt.Sprintf("http://localhost:%d/db/execute?pretty&timings", table.port)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return false, err
	}
	_, err = runQuery(req)
	// Handles resonpse.
	if err != nil {
		return false, err
	}
	return true, nil
}

// CreateTable creates a table to similate a read/write register with one row initialized to 0.
// Returns true if query succeeds.
func (table Table) CreateTable() (bool, error) {
	// transcation block to create new table an initiize with a default value of 0.
	create_query := fmt.Sprintf("\"CREATE TABLE %s (value INTEGER)\"", table.name)
	initialize_query := fmt.Sprintf("\"INSERT INTO %s(value) VALUES(0)\"", table.name)
	query := "[" + create_query + ", " + initialize_query + "]"
	endpoint := fmt.Sprintf("http://localhost:%d/db/execute?pretty&timings&transaction", table.port)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return false, err
	}
	_, err = runQuery(req)
	// Handles resonpse.
	if err != nil {
		return false, err
	}
	return true, nil
}

func (table Table) DeleteTable() (bool, error) {
	query := fmt.Sprintf("[\"DROP TABLE %s\"]", table.name)
	endpoint := fmt.Sprintf("http://localhost:%d/db/execute?pretty&timings", table.port)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return false, err
	}
	_, err = runQuery(req)
	// Handles resonpse.
	if err != nil {
		return false, err
	}
	return true, nil

}

func runQuery(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return resp, errors.New("Query did not succeed with status " + resp.Status)
	}

	return resp, nil
}

func main() {
	t := NewTable(4001, "test")
	fmt.Println(t.CreateTable())
	fmt.Println(t.Write(1))
	// fmt.Println(t.Write(2))
	fmt.Println(t.Read())
	fmt.Println(t.DeleteTable())

}
