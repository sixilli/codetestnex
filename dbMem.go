package main

import (
	"fmt"
    "errors"
    "sync"
)

// Person data structure that will serve as a schema
type Person struct {
	FirstName string
	LastName  string 
	Age       int    
}

// DBMem to solve race conditions will be managed by the JobStore
type DBMem struct {
    data     map[int]Person
    history  *HistoryStore
    sync.RWMutex
}

// InitDBM - Initialize DB in memory that will serve as the parent interface
func InitDBM() DBMem {
	return DBMem{
        data:    make(map[int]Person),
        history: InitHistoryStore(),
    }
}

// Get - Return db contents as a slice of structs
func (m *DBMem) Get() []Person {
    m.RLock()
    defer m.RUnlock()

	if len(m.data) == 0 {
		fmt.Println("Database is empty")
		return []Person{}
	}

	var output []Person
    for key := range m.data {
        output = append(output, m.data[key])
    }
	return output
}

// PrintDB - print contents of the DB with IDs
func (m *DBMem) PrintDB() {
    m.RLock()
    defer m.RUnlock()

    for i := 0; i < len(m.data); i++ {
        fmt.Println("ID:", i, m.data[i])
    }
}

// Read - returns list of requested ids
func (m *DBMem) Read(id int) (Person, error) {
    m.RLock()
    defer m.RUnlock()

    if len(m.data) < id {
		fmt.Println("ID is out of range")
		return Person{}, errors.New("ID is out of range")
	}
    
    v, ok := m.data[id]
    if !ok {
        return Person{}, errors.New("ID not found")
    }

    return v, nil
}

// Insert new row into the db
func (m *DBMem) Insert(data Person) {
    m.Lock()
    defer m.Unlock()

	id := len(m.data)
    m.data[id] = data
    m.history.Append("INSERT", id, data)
}

// Update row
// Could be greedy, but currently searching for proper ID
func (m *DBMem) Update(idToUpdate int, data Person) {
    m.Lock()
    defer m.Unlock()

	if len(m.data) <= idToUpdate {
		fmt.Println("ID is out of range")
		return
	}
    m.data[idToUpdate] = data
    m.history.Append("UPDATE", idToUpdate, data)
}
	
// Delete row with ID and reindex
func (m *DBMem) Delete(idToDelete int) {
    m.Lock()
    defer m.Unlock()

	if len(m.data) <= idToDelete {
		fmt.Println("ID", idToDelete,"is out of range")
		return
	}
    entryToDelete := m.data[idToDelete]
    delete(m.data, idToDelete)

    // Reindex database where ID > deleted
    // I think with using .Lock() m.data wasn't updating immediatly so I ignore the key
    tempMap := make(map[int]Person)
    for k, v := range m.data {
        switch{
        case k > idToDelete:
            tempMap[k-1] = v
        case k == idToDelete:
        default:
            tempMap[k] = v
        }
    }

    m.data = tempMap
    m.history.Append("DELETE", idToDelete, entryToDelete)
}
