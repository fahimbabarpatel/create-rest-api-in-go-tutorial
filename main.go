// main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Article - Our struct for all articles

type CurrentStatus struct {
	_id      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	VRN      string             `json:"vrn,omitempty" bson:"vrn,omitempty"`
	DriverId int                `json:"driverid,omitempty" bson:"driverid,omitempty"`
}

var client *mongo.Client

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/currentstatus", createCS).Methods("POST")
	myRouter.HandleFunc("/currentstatus/{id}", getCS).Methods("GET")
	myRouter.HandleFunc("/currentstatus", GetAll).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", myRouter))
}

func createCS(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")
	var currentstatus CurrentStatus
	json.NewDecoder(r.Body).Decode(&currentstatus)

	collection := client.Database("TelemetryDB_2").Collection("CurrentStatus")

	updateResult, err := collection.InsertOne(context.TODO(), currentstatus)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updates a single document: %+v\n", updateResult)
	json.NewEncoder(w).Encode(updateResult)

}

func getCS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var currentstatus CurrentStatus
	collection := client.Database("TelemetryDB_2").Collection("CurrentStatus")

	err := collection.FindOne(context.TODO(), CurrentStatus{_id: id}).Decode(&currentstatus)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(currentstatus)

}

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var currentstatuss []CurrentStatus
	collection := client.Database("TelemetryDB_2").Collection("CurrentStatus")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var currentstatus CurrentStatus
		cursor.Decode(&currentstatus)
		currentstatuss = append(currentstatuss, currentstatus)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(currentstatuss)
}
func main() {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:26552")
	client, _ = mongo.Connect(context.TODO(), clientOptions)
	handleRequests()
}
