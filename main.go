package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting")

    // Test data
    person := Person{
        FirstName: "Alec",
        LastName: "P",
        Age:  6,
    }

    // Test data
    newPerson := Person{
        FirstName: "Stan",
        LastName: "Guy",
        Age:  10,
    }

	// Initializing in memory db
	dm := InitDBM()
    for i := 0; i < 20;  i++ {
        dm.Insert(person)
    }
    dm.PrintDB()
	dm.Update(5, newPerson)
	dm.Delete(10)
    dm.PrintDB()


	// Uses buntdb
	InitDB()
	TestAPIGetAll()
	CreateRoutes()
}
