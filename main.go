package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting")

	// Uses buntdb
	InitDB()

	// In memory version
	dm := InitDBM()

    person := Person{
        FirstName: "Alec",
        LastName: "P",
        Age:  6,
    }

    newPerson := Person{
        FirstName: "Stan",
        LastName: "Guy",
        Age:  10,
    }

	dm.Insert(person)
	dm.Insert(person)
	dm.Insert(person)
	fmt.Println(dm.Get(), "\n")
	dm.Update(2, newPerson)
	dm.Delete(3)
	fmt.Println(dm.Get())

	TestAPIGetAll()
	CreateRoutes()
}
