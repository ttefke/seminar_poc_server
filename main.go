package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type LuxMessage struct {
	DeviceID  int64 `json:"device_id"`
	MessageID int64 `json:"message_id"`
	Lux       int64 `json:"lux"`
	Time      int64 `json:"unixTimestamp"`
}

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
	fmt.Printf("%s\n", m)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/lux", luxHandler).Methods("POST")

	fmt.Println("Starting HTTP listener")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatalln("HTTP listener failed: ", err)
	}
}
