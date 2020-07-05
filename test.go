package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type testCase struct {
	input    string
	expected string
}

// Could add a higher level function to run all tests

// TestAPIGetAll tests that the dummy data works as expected.
func TestAPIGetAll() error {
	testRead := testCase{
		input:    "",
		expected: `[{"FirstName":"Alec", "LastName":"Perro", "Age":5}]`,
	}

    query, err := dm.Read(1)
    if err != nil {
        log.Fatal(err)
    }

	jsonify, err := json.Marshal(query)
	if err != nil {
		log.Fatal(err)
	}

	if testRead.expected != string(jsonify) {
		return errors.New("testDB failed")
	}

	fmt.Println("Tests passed")
	return nil
}
