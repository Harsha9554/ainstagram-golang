package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name 	string `json:"name"`
	Email 	string `json:"email"`
	Password string `json:"password"`
}
var Users []User

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}


func usersEndpoint(w http.ResponseWriter, r *http.Request){
	var userCollection = db().Database("ainstagram").Collection("users")
	if r.URL.Path == "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "POST":
		w.Header().Set("content-type", "application/json")
		var user User
		fmt.Println(user)
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			fmt.Print(err)
		}
		insertResult, err:=
		userCollection.InsertOne(context.TODO(), user)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(insertResult.InsertedID)
	case "GET":
		w.Header().Set("content-type", "application/json")
		var users []User
		cursor, err := userCollection.Find(context.TODO(),bson.M{})
		if err != nil {
			fmt.Print(err)
		}
		defer cursor.Close(context.TODO())
		for cursor.Next(context.TODO()) {
			var user User
			cursor.Decode(&user)
			users = append(users, user)
		}
		
		if err := cursor.Err(); err != nil {
			fmt.Print(err)
		}
		json.NewEncoder(w).Encode(users)
	}

}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/users", usersEndpoint)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func db() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb+srv://harsha_9554:qwertyuiop@cluster0.94brn.mongodb.net/ainstagram?retryWrites=true&w=majority")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(),nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("mongodb connected")
	return client
}

func main() {
	handleRequests()
}