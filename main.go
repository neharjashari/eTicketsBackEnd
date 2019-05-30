package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	DateCreated string `json:"date_created"`
	Content     string `json:"content"`
	Photo       string `json:"photo"`
}

type User struct {
	Token string `json:"token"`
	Event Event  `json:"event"`
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/event/{token}", createEventHandler).Methods("POST")
	router.HandleFunc("/event/{token}", getEventsHandler).Methods("GET")
	router.HandleFunc("/events/{token}/{id}", getEventHandler).Methods("GET")

	//// Set http to listen and serve for different requests in the port found in the GetPort() function
	//err := http.ListenAndServe(GetPort(), router)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}

	log.Fatal(http.ListenAndServe(":8080", router))
}

func createEventHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	urlVars := mux.Vars(r)

	userEvent := User{}
	event := Event{}

	var error = json.NewDecoder(r.Body).Decode(event)
	if error != nil {
		fmt.Println("Error made: ", error)
		return
	}

	log.Print(event)

	token := urlVars["token"]

	userEvent.Event = event
	userEvent.Token = token

	// Connecting to DB
	client := mongoConnect()

	// Specifying the specific collection which is going to be used
	collection := client.Database("etickets").Collection("events")

	// Checking for duplicates
	duplicate := checkForDuplicates(event.Title, token)

	// If there are not duplicates Insert that track to the collection
	if !duplicate {

		res, err := collection.InsertOne(context.Background(), userEvent)
		if err != nil {
			log.Fatal(err)
		}
		id := res.InsertedID

		if id == nil {
			http.Error(w, "", 300)
		}

		// Encoding the ID of the track that was just added to DB
		err = json.NewEncoder(w).Encode(event.ID)
		if err != nil {
			fmt.Println("Error made while encoding with JSON, : ", err)
			return
		}

	} else {

		// Notifying the user that the IGC File posted is already in our DB
		http.Error(w, "409 Conflict - The Event you entered is already in our database!", http.StatusConflict)
		return
	}

}

func getEventsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Using mux router to save the ID variable defined in the requested path
	urlVars := mux.Vars(r)
	token := urlVars["token"]

	client := mongoConnect()

	events := getAllEvents(client, token)

	// Encoding all IDs of the track in IgcFilesDB
	err := json.NewEncoder(w).Encode(events)
	if err != nil {
		fmt.Println("Error made while encoding with JSON, : ", err)
		return
	}

}

func getEventHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Using mux router to save the ID variable defined in the requested path
	urlVars := mux.Vars(r)
	token := urlVars["token"]

	client := mongoConnect()

	events := getAllEvents(client, token)

	event := Event{}

	for i := range events {
		if events[i].ID == urlVars["id"] {
			event = events[i]

			err := json.NewEncoder(w).Encode(event)
			if err != nil {
				fmt.Println("Error made while encoding with JSON, : ", err)
				return
			}
		}
	}

}
