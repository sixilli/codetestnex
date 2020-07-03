package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateRoutes uses gorilla mux router to handle requests
// API update example
// localhost:8080/update/1/bill/gates/30
func CreateRoutes() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homePage).Methods("GET")
	r.HandleFunc("/create/{firstName}/{lastName}/{age}", createEntry).Methods("POST")
	r.HandleFunc("/read", readEntry).Methods("GET")
	r.HandleFunc("/delete/{id}", deleteEntry).Methods("DELETE")
	r.HandleFunc("/update/{id}/{firstName}/{lastName}/{age}", updateEntry).Methods("PUT")
	r.HandleFunc("/live", ServeLive)
	r.HandleFunc("/ws", ServeWs)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

// All other route helper functions handle api calls
func createEntry(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creating entry \n")

	vars := mux.Vars(r)
	firstName := vars["firstName"]
	lastName := vars["lastName"]
	age, err := strconv.Atoi(vars["age"])
	if err != nil {
		log.Fatal(err)
	}

	newEntry := Person{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
	APINewEntry(newEntry)
}

// Returns an array of JSON objects
func readEntry(w http.ResponseWriter, r *http.Request) {
	get := APIGetAll()
	json.NewEncoder(w).Encode(get)
}

func deleteEntry(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Deleting Entry")

	vars := mux.Vars(r)
	id := vars["id"]

	APIDeleteEntry(id)
}

func updateEntry(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Updating")

	vars := mux.Vars(r)
	id := vars["id"]
	firstName := vars["firstName"]
	lastName := vars["lastName"]
	age, err := strconv.Atoi(vars["age"])
	if err != nil {
		log.Fatal(err)
	}

	toUpdate := Person{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}

	APIUpdateEntry(id, toUpdate)
}
