package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// *** DB METHODS *** //


// This function connects the API with Mongo Database and returns that connection
func mongoConnect() *mongo.Client {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	//clientOptions := options.Client().ApplyURI("mongodb://neharjashari:nerkoid17051998@ds263856.mlab.com:63856/?authSource=etickets")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Connected to MongoDB!")

	return client
}


// Get all events with token
func getAllEvents(client *mongo.Client, token string) []Event {
	db := client.Database("etickets")
	collection := db.Collection("events")

	cursor, err := collection.Find(context.Background(), bson.D{{"token", token}})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	userEvent := User{}
	events := []Event{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&userEvent)
		if err != nil {
			log.Fatal(err)
		}
		events = append(events, userEvent.Event)
	}

	return events
}


// Get all tickets with token
func getAllTickets(client *mongo.Client, token string) []Event {
	db := client.Database("etickets")
	collection := db.Collection("tickets")

	cursor, err := collection.Find(context.Background(), bson.D{{"token", token}})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	userEvent := User{}
	events := []Event{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&userEvent)
		if err != nil {
			log.Fatal(err)
		}
		events = append(events, userEvent.Event)
	}

	return events
}


// Get all tracks
func getAllEventsForMainActivity(client *mongo.Client) []Event {
	db := client.Database("etickets")
	collection := db.Collection("events")

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	userEvent := User{}
	all := []Event{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&userEvent)
		if err != nil {
			log.Fatal(err)
		}
		all = append(all, userEvent.Event)
	}

	return all
}


func checkForDuplicates(client *mongo.Client, ID string, Title string, token string) (bool, string) {

	events := getAllEvents(client, token)
	bool := false
	response := ""

	for i := range events {
		if events[i].ID == ID {
			bool = true
			response = "An event with that ID already is in your database."
		} else if events[i].Title == Title {
			bool = true
			response = "An event with that Title already is in your database."
		}
	}

	return bool, response
}


func checkForTicketDuplicates(client *mongo.Client, ID string, Title string, token string) (bool, string) {

	events := getAllTickets(client, token)
	bool := false
	response := ""

	for i := range events {
		if events[i].ID == ID {
			bool = true
			response = "An event with that ID already is in your database."
		} else if events[i].Title == Title {
			bool = true
			response = "An event with that Title already is in your database."
		}
	}

	return bool, response
}


// Delete all events
func deleteAllEvents(client *mongo.Client) {
	db := client.Database("etickets")
	collection := db.Collection("events")

	collection.DeleteMany(context.Background(), bson.D{})
}

// Delete all events
func deleteAllTickets(client *mongo.Client) {
	db := client.Database("etickets")
	collection := db.Collection("tickets")

	collection.DeleteMany(context.Background(), bson.D{})
}


// Get all the data
func getAllEventsAdmin(client *mongo.Client) []User {
	db := client.Database("etickets")
	collection := db.Collection("events")

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	userEvent := User{}
	all := []User{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&userEvent)
		if err != nil {
			log.Fatal(err)
		}
		all = append(all, userEvent)
	}

	return all
}
