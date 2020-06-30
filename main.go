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
	dm.Insert()
	dm.Insert()
	fmt.Println(dm.Get())

	TestAPIGetAll()
	CreateRoutes()
}
