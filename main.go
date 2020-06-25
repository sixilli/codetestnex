package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting")
	InitDB()
	TestAPIGetAll()
	CreateRoutes()
}
