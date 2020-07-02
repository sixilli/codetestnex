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

	// Go through each row and jsonify the struct
	var jsonOutput []string
	for i := 0; i < len(m.Rows); i++ {
		jsonOutput = append(jsonOutput, m.Rows[i].Data)
	}

	return jsonOutput
}

// Insert new row into the db
func (m *DBMem) Insert(data Person) {
	id := len(m.Rows) + 1

    bytes, err := json.Marshal(data)
    if err != nil {
        fmt.Println(err)
    }

	newRow := dbRow{ID: id, Data: string(bytes)}
	m.Rows = append(m.Rows, newRow)
}

// Update row
func (m *DBMem) Update(idToUpdate int, data Person) {
	if len(m.Rows) < idToUpdate {
		fmt.Println("ID is out of range")
		return
	}

    bytes, err := json.Marshal(data)
    if err != nil {
        fmt.Println(err)
    }

	for i := 0; i < len(m.Rows); i++ {
		if m.Rows[i].ID == idToUpdate {
			m.Rows[i].Data = string(bytes)
			return
		}
	}

	fmt.Println("ID not found")
}
	
// Delete row with ID
func (m *DBMem) Delete(idToDelete int) {
	if len(m.Rows) < idToDelete {
		fmt.Println("ID is out of range")
		return
	}

	if len(m.Rows) == idToDelete {
		m.Rows = m.Rows[:len(m.Rows)-1]
		return
	}
	m.Rows = append(m.Rows[:idToDelete], m.Rows[idToDelete+1:]...)
    m.reindexDb()
}

// This can be made faster
func (m *DBMem) reindexDb() {
    for i := 1; i < len(m.Rows); i++ {
        m.Rows[i].ID = i
    }
}
