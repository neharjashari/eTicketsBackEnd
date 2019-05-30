package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// *** DB METHODS *** //

// This function connects the API with Mongo Database and returns that connection
func mongoConnect() *mongo.Client {
	// Connect to MongoDB

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return client
}

// Get all tracks
func getAllEvents(client *mongo.Client, token string) []Event {
	db := client.Database("etickets")
	collection := db.Collection("events")

	cursor, err := collection.Find(context.Background(), bson.D{{"token", token}})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	events := []Event{}
	userEvent := User{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&userEvent)
		if err != nil {
			log.Fatal(err)
		}
		events = append(events, userEvent.Event)
	}

	return events
}

func checkForDuplicates(ID string, token string) bool {

	client := mongoConnect()

	events := getAllEvents(client, token)

	for i := range events {
		if events[i].ID == ID {
			return true
		}
	}

	return false
}

//// Delete all tracks
//func deleteAllTracks(client *mongo.Client) {
//	db := client.Database("etickets")
//	collection := db.Collection("events")
//
//	collection.DeleteMany(context.Background(), bson.NewDocument())
//}

//// Count all tracks
//func countAllTracks(client *mongo.Client) int64 {
//	db := client.Database("igcfiles")
//	collection := db.Collection("track")
//
//	// Count the tracks
//	count, _ := collection.Count(context.Background(), nil, nil)
//
//	return count
//}
