package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"regexp"
	// "math/rand"
	// "strconv"
)

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	DateCreated string `json:"date_created"`
	Content     string `json:"content"`
	Price       string `json:"price"`
}

type User struct {
	Token string `json:"token"`
	Event Event  `json:"event"`
}

type MetaInformation struct {
	Info    []string `json:"info"`
	Version string   `json:"version"`
}

 //// Get the Port from the environment so we can run on Heroku
 //func GetPort() string {
 //	var port = os.Getenv("PORT")
 //	// Set a default port if there is nothing in the environment
 //	if port == "" {
 //		port = "4747"
 //		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
 //	}
 //	return ":" + port
 //}


func main() {

	router := mux.NewRouter()

	router.HandleFunc("/info", infoHandler).Methods("GET")
	router.HandleFunc("/events", getAllEventsHandler).Methods("GET")
	router.HandleFunc("/events/{id}", getEventWithIdHandler).Methods("GET")
	router.HandleFunc("/event/{token}", createEventHandler).Methods("POST")
	router.HandleFunc("/event/{token}/tickets", getTicketsHandler).Methods("GET")
	router.HandleFunc("/event/{token}/tickets", createTicketsHandler).Methods("POST")
	router.HandleFunc("/event/{token}", getEventsHandler).Methods("GET")
	router.HandleFunc("/event/{token}/{id}", getEventHandler).Methods("GET")

	router.HandleFunc("/admin", adminHandler)
	router.HandleFunc("/admin/tickets", adminTicketHandler)
	router.HandleFunc("/admin/event/{id}", deleteEventHandler)

	//// Set http to listen and serve for different requests in the port found in the GetPort() function
	//err := http.ListenAndServe(GetPort(), router)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}

	log.Fatal(http.ListenAndServe(":8000", router))
}

func infoHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "400 - Bad Request, Wrong method", http.StatusBadRequest)
		return
	}

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

	if r.Method != http.MethodPost {
		http.Error(w, "400 - Bad Request, Wrong method", http.StatusBadRequest)
		return
	}

	urlVars := mux.Vars(r)

	userEvent := User{}
	event := Event{}

	var error = json.NewDecoder(r.Body).Decode(&event)
	if error != nil {
		fmt.Println("Error made: ", error)
		return
	}

	token := urlVars["token"]

	// Validate token
	valid := validateToken(token)
	if !valid {
		http.Error(w, "400 Bad Request - The token you wrote is not valid.", http.StatusBadRequest)
		return
	}

	userEvent.Token = token
	userEvent.Event = event

	// Connecting to DB
	client := mongoConnect()

	// Specifying the specific collection which is going to be used
	collection := client.Database("etickets").Collection("events")

	//Checking for duplicates
	duplicate, response := checkForDuplicates(client, event.ID, event.Title, token)

	if duplicate {

		http.Error(w, "409 Conflict - " + response, http.StatusConflict)
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

	if r.Method != http.MethodGet {
		http.Error(w, "400 - Bad Request, Wrong method", http.StatusBadRequest)
		return
	}

	urlVars := mux.Vars(r)
	token := urlVars["token"]

	// Validate token
	valid := validateToken(token)
	if !valid {
		http.Error(w, "400 Bad Request - The token you wrote is not valid.", http.StatusBadRequest)
		return
	}

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

	if r.Method != http.MethodGet {
		http.Error(w, "400 - Bad Request, Wrong method", http.StatusBadRequest)
		return
	}

	urlVars := mux.Vars(r)
	token := urlVars["token"]

	// Validate token
	valid := validateToken(token)
	if !valid {
		http.Error(w, "400 Bad Request - The token you wrote is not valid.", http.StatusBadRequest)
		return
	}

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


func getAllEventsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "400 - Bad Request, Wrong method", http.StatusBadRequest)
		return
	}

	client := mongoConnect()

	events := getAllEventsForMainActivity(client)

	err := json.NewEncoder(w).Encode(events)
	if err != nil {
		fmt.Println("Error made while encoding with JSON, : ", err)
		return
	}
}


func getEventWithIdHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "400 - Bad Request, Wrong method", http.StatusBadRequest)
		return
	}

	urlVars := mux.Vars(r)

	client := mongoConnect()

	events := getAllEventsForMainActivity(client)

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


func getTicketsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "400 - Bad Request, Wrong method", http.StatusBadRequest)
		return
	}

	urlVars := mux.Vars(r)
	token := urlVars["token"]

	// Validate token
	valid := validateToken(token)
	if !valid {
		http.Error(w, "400 Bad Request - The token you wrote is not valid.", http.StatusBadRequest)
		return
	}

	client := mongoConnect()

	events := getAllTickets(client, token)

	err := json.NewEncoder(w).Encode(events)
	if err != nil {
		fmt.Println("Error made while encoding with JSON, : ", err)
		return
	}
}


func createTicketsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "400 - Bad Request, Wrong method", http.StatusBadRequest)
		return
	}

	urlVars := mux.Vars(r)

	userEvent := User{}
	event := Event{}

	var error = json.NewDecoder(r.Body).Decode(&event)
	if error != nil {
		fmt.Println("Error made: ", error)
		return
	}

	token := urlVars["token"]

	// Validate token
	valid := validateToken(token)
	if !valid {
		http.Error(w, "400 Bad Request - The token you wrote is not valid.", http.StatusBadRequest)
		return
	}

	userEvent.Token = token
	userEvent.Event = event

	// Connecting to DB
	client := mongoConnect()

	// Specifying the specific collection which is going to be used
	collection := client.Database("etickets").Collection("tickets")

	//Checking for duplicates
	duplicate, response := checkForTicketDuplicates(client, event.ID, event.Title, token)

	if duplicate {

		http.Error(w, "409 Conflict - " + response, http.StatusConflict)
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


//-----------------------------------------------------------------------
// ADMIN
//-----------------------------------------------------------------------
func adminHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		w.Header().Set("Content-Type", "application/json")

		client := mongoConnect()

		all := getAllEventsAdmin(client)

		err := json.NewEncoder(w).Encode(all)
		if err != nil {
			fmt.Println("Error made while encoding with JSON, : ", err)
			return
		}

	case "DELETE":
		client := mongoConnect()

		deleteAllEvents(client)

		response := "All events are deleted!"

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			fmt.Println("Error made while encoding with JSON, : ", err)
			return
		}

	default:
		http.Error(w, "400 - Bad Request, Wrong method", http.StatusBadRequest)
		return

	}
}


func adminTicketHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, "400 - Bad Request, Wrong method", http.StatusBadRequest)
		return
	}

	client := mongoConnect()

	deleteAllTickets(client)

	response := "All tickets are deleted!"

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Error made while encoding with JSON, : ", err)
		return
	}
	
}


func deleteEventHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, "400 - Bad Request, Wrong method", http.StatusBadRequest)
		return
	}

	urlVars := mux.Vars(r)

	client := mongoConnect()

	id := urlVars["id"]

	deleteEvent(client, id)

	response := "The event with ID: " + id + " , has been deleted!"

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Error made while encoding with JSON, : ", err)
		return
	}
	
}



func validateToken(token string) bool {

	valid := false

	// Regular Expression to check for Token validity, the ID can only be with the same format as UUID
	regExToken, _ := regexp.Compile("/^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/")

	if !regExToken.MatchString(token) {
		valid = true
	}

	return valid
}


// func randomID() string {
// 	randomNumber := rand.Intn(1000)
// 	strRandom := strconv.Itoa(randomNumber)
// 	return strRandom
// }