package main

import (
	"fmt"
)

var dm = InitDBM()

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
    for i := 0; i < 20;  i++ {
        dm.Insert(person)
    }
    dm.PrintDB()
    fmt.Println("--------------Testing Update and Delete--------------")
	dm.Update(5, newPerson)
	dm.Delete(10)
	dm.Delete(2)
	dm.Delete(4)
	dm.Delete(19)
    dm.PrintDB()


	// Uses buntdb
	InitDB()
	TestAPIGetAll()
	CreateRoutes()
}
