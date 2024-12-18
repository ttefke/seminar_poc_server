package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type LuxMessage struct {
	DeviceID  int64  `json:"device_id"`
	MessageID int64  `json:"message_id"`
	Lux       int64  `json:"lux"`
	Time      int64  `json:"unixTimestamp"`
	UUID      string `json:"uuid"`
}

var (
	collection *mongo.Collection
)

func luxHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error occurred while reading request: ", err)
	}

	body = []byte(body)

	if !json.Valid(body) {
		fmt.Println("Invalid JSON received")
		return
	}

	var m LuxMessage
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Println("Error occurred while unmarshalling JSON: ", err)
	}
	m.Time = time.Now().Unix()
	m.UUID = uuid.New().String()

	result, err := collection.InsertOne(context.TODO(), m)
	fmt.Printf("Inserted document with id %v\n", result.InsertedID)
}

func main() {
	fmt.Println("Connecting to DB")

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	hostname := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")

	if len(username) == 0 || len(password) == 0 ||
		len(hostname) == 0 || len(dbname) == 0 {
		fmt.Println("No database credentials provided")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+username+":"+password+"@"+hostname+":27017/"+dbname))
	if err != nil {
		log.Fatalln("Could not connect to DB:", err)
		os.Exit(1)
	}
	collection = client.Database(dbname).Collection("light")

	r := mux.NewRouter()
	r.HandleFunc("/lux", luxHandler).Methods("POST")

	fmt.Println("Starting HTTP listener")
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatalln("HTTP listener failed: ", err)
	}
	defer cancel()
	defer client.Disconnect(ctx)
}
