package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	configFilename = "config.json"
)

// Configuration describes the configuration format
type Configuration struct {
	Server Server `json:"server"`
}

// Server describes the server configuration
type Server struct {
	ServerPort int32  `json:"port"`
	Endpoint   string `json:"endpoint"`
}

func main() {
	file, err := os.Open(configFilename)
	if err != nil {
		log.Fatalf("Could not load file: %s", configFilename)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	config := Configuration{}
	err = decoder.Decode(config)
	if err != nil {
		log.Fatalf("Could not decode config file: %v", err)
	}

	http.HandleFunc(config.Server.Endpoint, handle)

	port := strconv.FormatInt(int64(config.Server.ServerPort), 10)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	event := r.Header.Get("X-GitHub-Event")
	body, err := r.GetBody()
	if err != nil {
		log.Printf("Could not extract body of incomming message")
		return
	}

	switch event {
	case "push":
		handlePush(body)
	}
}

func handlePush(body io.ReadCloser) (map[string]interface{}, error) {
	var data []byte
	body.Read(data)

	var payload map[string]interface{}
	err := json.Unmarshal(data, payload)

	if err != nil {
		log.Printf("Could not unmarshal json payload: %v", err)
		return nil, err
	}

	return payload, nil
}
