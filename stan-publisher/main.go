package main

import (
	"log"
	"os"

	"github.com/nats-io/stan.go"
)

func main() {
	jsonFile := "./stan-publisher/model.json"

	inputJSON, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf("Can`t read file: %v", err)
	}

	sc, err := stan.Connect("test-cluster", "pub-client")
	if err != nil {
		log.Fatalf("Can't connect: %v.", err)
	}

	defer sc.Close()

	err = sc.Publish("orders", inputJSON)
	if err != nil {
		log.Printf("Something wrong when trying publish:%v", err)
	}
	log.Print("data successfully sent to the channel")
}
