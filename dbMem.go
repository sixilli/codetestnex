package main

import (
	"encoding/json"
	"fmt"
)

type dbRow struct {
	ID   int    `json:"ID"`
	Data string `json:"Data"`
}

// DBMem - This is my take on learning how to create a nice interface
// to the following group of functions.
type DBMem struct {
	Rows []dbRow
}

func InitDBM() DBMem {
	return DBMem{}
}

// Get - Return db contents as string
func (m *DBMem) Get() []string {
	if len(m.Rows) == 0 {
		fmt.Println("Database is empty")
		return []string{}
	}
	out, err := json.Marshal(m.Rows)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	var jsonOutput []string
	for i := 0; i < len(out); i++ {
		jsonOutput = append(jsonOutput, string(out[i]))
	}
	return jsonOutput
}

// Insert new row into the db
func (m *DBMem) Insert() {
	id := len(m.Rows) + 1
	newRow := dbRow{ID: id, Data: `{FistName: "Alec", LastName: "P" Age: 4, }`}
	m.Rows = append(m.Rows, newRow)
}
