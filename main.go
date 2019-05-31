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

type MetaInformation struct {
	Info    []string `json:"info"`
	Version string   `json:"version"`
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/info", infoHandler).Methods("GET")
	router.HandleFunc("/event/{token}", createEventHandler).Methods("POST")
	router.HandleFunc("/event/{token}", getEventsHandler).Methods("GET")
	router.HandleFunc("/event/{token}/{id}", getEventHandler).Methods("GET")

	router.HandleFunc("/admin", deleteAllAdmin).Methods("DELETE")
	router.HandleFunc("/admin", getAllAdmin).Methods("GET")

	//// Set http to listen and serve for different requests in the port found in the GetPort() function
	//err := http.ListenAndServe(GetPort(), router)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}

	log.Fatal(http.ListenAndServe(":8000", router))
}

func infoHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	metaInfo := &MetaInformation{}

	metaInfo.Info = append(metaInfo.Info, "E-TICKETS Application")
	metaInfo.Info = append(metaInfo.Info, "Back End Service")
	metaInfo.Info = append(metaInfo.Info, "Service for Creating Events")
	metaInfo.Version = "v2.2.0"

	err := json.NewEncoder(w).Encode(metaInfo)
	if err != nil {
		fmt.Println("Error made while encoding with JSON, : ", err)
		return
	}

}

func createEventHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	urlVars := mux.Vars(r)

	userEvent := User{}
	event := Event{}

	var error = json.NewDecoder(r.Body).Decode(&event)
	if error != nil {
		fmt.Println("Error made: ", error)
		return
	}

	token := urlVars["token"]

	userEvent.Token = token
	userEvent.Event = event

	// Connecting to DB
	client := mongoConnect()

	// Specifying the specific collection which is going to be used
	collection := client.Database("etickets").Collection("events")

	//Checking for duplicates
	duplicate := checkForDuplicates(client, event.Title, token)

	if duplicate {

		http.Error(w, "409 Conflict - The Event you entered is already in our database!", http.StatusConflict)
		return

	} else {

		res, err := collection.InsertOne(context.Background(), userEvent)
		if err != nil {
			log.Fatal(err)
			log.Print("Error")
		}

		id := res.InsertedID

		if id == nil {
			http.Error(w, "", 300)
		}

		log.Print(id)

		err = json.NewEncoder(w).Encode(event.ID)
		if err != nil {
			fmt.Println("Error made while encoding with JSON, : ", err)
			return
		}
	}
}

func getEventsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	urlVars := mux.Vars(r)
	token := urlVars["token"]

	client := mongoConnect()

	events := getAllEvents(client, token)

	err := json.NewEncoder(w).Encode(events)
	if err != nil {
		fmt.Println("Error made while encoding with JSON, : ", err)
		return
	}
}

func getEventHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

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

func deleteAllAdmin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	client := mongoConnect()

	deleteAllEvents(client)

	response := "All events are deleted!"

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Error made while encoding with JSON, : ", err)
		return
	}
}

func getAllAdmin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	client := mongoConnect()

	getAllEventsAdmin(client)

	response := "All events are deleted!"

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Error made while encoding with JSON, : ", err)
		return
	}
}
