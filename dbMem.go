package main

import (
	"encoding/json"
	"fmt"
	"strings"
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

	// Go through each row and jsonify the struct
	var jsonOutput []string
	for i := 0; i < len(m.Rows); i++ {
		row := &m.Rows[i]
		b, err := json.Marshal(row)
		if err != nil {
			fmt.Println(err)
		}

		// Kinda hacky, but the JSON object had added in escape characters
		jsonOutput = append(jsonOutput, strings.ReplaceAll(string(b), `\`, ""))
	}

	return jsonOutput
}

// Insert new row into the db
func (m *DBMem) Insert() {
	id := len(m.Rows) + 1
	newRow := dbRow{ID: id, Data: `{"FistName": "Alec", "LastName": "P", "Age": 4}`}
	m.Rows = append(m.Rows, newRow)
}
