package main

import (
	"fmt"
)

var dm DBMem

func main() {
	fmt.Println("Starting")

    // Test data
    newPerson := Person{
        FirstName: "Stan",
        LastName: "Guy",
        Age:  10,
    }

	// Initializing in memory db
    dm = InitDBM()
    for i := 0; i < 20;  i++ {
        dm.Insert(Person{
            FirstName: "Alec",
            LastName: "p",
            Age: i,
        })
    }
    dm.PrintDB()
    fmt.Println("--------------Testing Update and Delete--------------")
	dm.Update(5, newPerson)
	dm.Delete(10)
	dm.Delete(2)
	dm.Delete(4)
	dm.Delete(19) // This serves as a test for error handeling
    dm.PrintDB()

	CreateRoutes()
}
