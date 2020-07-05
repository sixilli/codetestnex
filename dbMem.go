package main

import (
	"encoding/json"
	"fmt"
    "time"
    "errors"
)

type dbEvent struct {
    Op        string `json:"Op"`
    ID        int    `json:"ID"`
    Timestamp string `json:"Time"`
    Data      string `json:"Data"`
}

type dbRow struct {
	ID   int    `json:"ID"`
	Data string `json:"Data"`
}

// DBMem - This is my take on learning how to create a nice interface
// to the following group of functions.
// Might be bad design to add history to this, but makes things easier
type DBMem struct {
	Rows []dbRow
    History []dbEvent
}

// InitDBM - Initialize DB in memory that will serve as the parent interface
func InitDBM() DBMem {
	return DBMem{}
}

// Get - Return db contents as a slice of structs
func (m *DBMem) Get() []Person {
	if len(m.Rows) == 0 {
		fmt.Println("Database is empty")
		return []Person{}
	}

	// Go through each row and convert json to a struct
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

// Read - returns list of requested ids
func (m *DBMem) Read(id int) (Person, error) {
    if len(m.Rows) < id {
		fmt.Println("ID is out of range")
		return Person{}, errors.New("ID is out of range")
	}

	for i := 0; i < len(m.Rows); i++ {
		if m.Rows[i].ID == id {
            personStruct := Person{}
            json.Unmarshal([]byte(m.Rows[i].Data), &personStruct)
			return personStruct, nil
		}
	}
    return Person{}, errors.New("ID not found")
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
    m.appendHistory("INSERT", id, newRow.Data, time.Now().String())
}

// Update row
// Could be greedy, but currently searching for proper ID
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
            m.appendHistory("Update", m.Rows[i].ID, m.Rows[i].Data, time.Now().String())
			return
		}
	}

	fmt.Println("ID not found")
}
	
// Delete row with ID
func (m *DBMem) Delete(idToDelete int) {
	if len(m.Rows) < idToDelete {
		fmt.Println("ID", idToDelete,"is out of range")
		return
	}

    rowData := m.Rows[idToDelete]

	if len(m.Rows) == idToDelete {
		m.Rows = m.Rows[:len(m.Rows)-1]
		return
	}
	m.Rows = append(m.Rows[:idToDelete], m.Rows[idToDelete+1:]...)
    m.appendHistory("DELETE", rowData.ID, rowData.Data, time.Now().String())
    m.reindexDb(idToDelete)
}

// Start at deleted entry then decrement all rows after
func (m *DBMem) reindexDb(deletedId int) {
    for i := deletedId; i < len(m.Rows); i++ {
        m.Rows[i].ID = m.Rows[i].ID - 1
    }
}

func (m *DBMem) appendHistory(operation string, id int, data string, timestamp string) {
    newEvent := dbEvent{
        Op: operation,
        ID: id,
        Data: data,
        Timestamp: timestamp,
    }
    m.History = append(m.History, newEvent)
}
