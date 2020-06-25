package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/tidwall/buntdb"
)

var db *buntdb.DB

// Person data structure that represents NoSQL data
type Person struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Age       int    `json:"Age"`
}

// Persons data structure where each entry is a row from the NoSQL db
type Persons []Person

// APIGetAll function to handle returning all rows in the db
func APIGetAll() Persons {
	// Initialize an array of Persons to append to.
	var getRequest Persons
	// db.View serves as a read only
	// tx.Ascend allows for iteration over the data
	db.View(func(tx *buntdb.Tx) error {
		tx.Ascend("LastName", func(key, value string) bool {
			// Jsonify output and append
			var jsonValue Person
			json.Unmarshal([]byte(value), &jsonValue)
			getRequest = append(getRequest, Person{
				FirstName: jsonValue.FirstName,
				LastName:  jsonValue.LastName,
				Age:       jsonValue.Age,
			})
			return true
		})
		return nil
	})
	return getRequest
}

// APINewEntry function to handle new rows
func APINewEntry(p Person) {
	newEntry, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	// Len serves as an auto increment for the db, would be best have
	// this run once and updated, rather than for every insert.
	var len int

	db.View(func(tx *buntdb.Tx) error {
		len, _ = tx.Len()
		len++
		return nil
	})

	// Adds new entry to the db if ID is not found
	db.Update(func(tx *buntdb.Tx) error {
		fmt.Println(string(newEntry))
		tx.Set(strconv.Itoa(len), string(newEntry), nil)
		return nil
	})
}

// APIDeleteEntry function to handle deletions using ID number
func APIDeleteEntry(index string) {
	db.Update(func(tx *buntdb.Tx) error {
		tx.Delete(index)
		return nil
	})
}

// APIUpdateEntry function handles updating using the ID
func APIUpdateEntry(index string, p Person) {
	toUpdate, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	updateEntry := string(toUpdate)

	db.Update(func(tx *buntdb.Tx) error {
		tx.Set(index, updateEntry, nil)
		return nil
	})
}

// InitDB function to connect the DB and setup for operations
func InitDB() {
	// Running the db in memory
	db, _ = buntdb.Open(":memory:")

	// Creating the JSON object for each row
	db.CreateIndex("FirstName", "*", buntdb.IndexJSON("FirstName"))
	db.CreateIndex("LastName", "*", buntdb.IndexJSON("LastName"))
	db.CreateIndex("Age", "*", buntdb.IndexJSON("Age"))

	// Adding some dummy data
	db.Update(func(tx *buntdb.Tx) error {
		tx.Set("1", `{"FirstName":"Alec", "LastName":"Perro", "Age":5}`, nil)
		tx.Set("2", `{"FirstName":"Al", "LastName":"Peterson", "Age":6}`, nil)
		return nil
	})
}
