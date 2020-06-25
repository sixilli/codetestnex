package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Dead code for websockets
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client Dead code for websockets
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

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
	r.HandleFunc("/live", websocketEndpoint).Methods("GET")
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

// Leftover websocket code
func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

// Leftover websocket code
func websocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected")

	if r.URL.Path != "/live" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not found", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "index.html")

	reader(ws)
}
