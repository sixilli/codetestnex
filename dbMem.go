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

// Get - Return db contents as a slice of structs
func (m *DBMem) Get() []Person {
	if len(m.Rows) == 0 {
		fmt.Println("Database is empty")
		return []Person{}
	}

	// Go through each row and jsonify the struct
	var output []Person
	for i := 0; i < len(m.Rows); i++ {
        personStruct := Person{}
        json.Unmarshal([]byte(m.Rows[i].Data), &personStruct)

		output = append(output, personStruct)
	}

	return output
}

// PrintDB - print contents of the DB with IDs
func (m *DBMem) PrintDB() {
    for i := 0; i < len(m.Rows); i++ {
        fmt.Println("ID:", m.Rows[i].ID, m.Rows[i].Data)
    }
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
    m.reindexDb(idToDelete)
}

// Start at deleted entry then decrement all rows after
func (m *DBMem) reindexDb(deletedId int) {
    for i := deletedId; i < len(m.Rows); i++ {
        fmt.Println("Decrementing", m.Rows[i].ID)
        m.Rows[i].ID = m.Rows[i].ID - 1
    }
}
