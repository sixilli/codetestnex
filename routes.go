package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
    "text/template"

	"github.com/gorilla/mux"
)

// CreateRoutes uses gorilla mux router to handle requests
// API update example
// localhost:8080/update/1/bill/gates/30
func CreateRoutes() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homePage).Methods("GET")
	r.HandleFunc("/create/{firstName}/{lastName}/{age}", createEntry).Methods("POST")
	r.HandleFunc("/read/{id}", readEntry).Methods("GET")
	r.HandleFunc("/delete/{id}", deleteEntry).Methods("DELETE")
	r.HandleFunc("/update/{id}/{firstName}/{lastName}/{age}", updateEntry).Methods("PUT")
	r.HandleFunc("/live", liveUpdate).Methods("GET")
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

    dm.Insert(newEntry)
}

// Returns an array of JSON objects
func readEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        log.Println("Bad request:", err)
		w.WriteHeader(400)
		return
    }

	get, err := dm.Read(id)
    if err != nil {
        log.Println("Bad request:", err)
		w.WriteHeader(400)
		return
    }

	json.NewEncoder(w).Encode(get)
}

func deleteEntry(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Deleting Entry")

	vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        log.Println("Bad request, invalid id.", err)
		w.WriteHeader(400)
		return
    }

    dm.Delete(id)
}

func updateEntry(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Updating")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
        log.Println("Bad request, invalid id.", err)
		w.WriteHeader(400)
		return
    }

	firstName := vars["firstName"]
	lastName := vars["lastName"]

	age, err := strconv.Atoi(vars["age"])
	if err != nil {
        log.Println("Bad request, invalid id.", err)
		w.WriteHeader(400)
		return
	}

	toUpdate := Person{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}

	dm.Update(id, toUpdate)
}


func liveUpdate(w http.ResponseWriter, r *http.Request) {
        template := template.New("template")
        // "doc" is the constant that holds the HTML content
        template.New("doc").Parse(liveHTML)
        context := Context{
            History: dm.History,
        }
        template.Lookup("doc").Execute(w, context)
}

const liveHTML = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <title>Nexoff Code Test - Live</title>
    </head>
    <body>
        <h2>Live Updates</h2>
        <ul>
            {{range .History}}
                <li>{{.}}</li>
            {{end}}
        </ul>
    </body>
</html>
`
