package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"service/service/internal/config"
	"service/service/internal/server"
	"service/service/pkg/cache"
	"service/service/pkg/database"
	"service/service/pkg/stansub"

	"github.com/nats-io/stan.go"
)

const configPath = "service/config/config.yaml"

func main() {

	// LOAD CONFIG
	config := config.LoadConfigYaml(configPath)
	log.Printf("Loaded Config: %v\n", config)

	// Connect TO STAN
	sc, err := stansub.ConnectToStan(config.STAN.ClusterID, config.STAN.ClientID)
	if err != nil {
		log.Fatalf("Can't connect to stan: %v.", err)
	}
	defer sc.Close()
	log.Printf("Connected to [%s] clusterID: [%s] clientID: [%s]\n", stan.DefaultNatsURL, config.STAN.ClusterID, config.STAN.ClientID)

	// CONNECT TO DB
	db, err := database.ConnectDB(&config.DB)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Printf("Connected to database-server: %v", db)

	// INIT CACHE
	csh, err := cache.New(0, 0, db)
	if err != nil {
		log.Fatalf("Can't create cache: %v", err)
	}
	log.Print("Cache succesfully created")

	// CREATE AND RUN SERVER
	srv := server.CreateServer(&config.HTTPServer, csh)
	go server.RunServer(srv)
	log.Printf("HTTP-server running on [%v] now", config.HTTPServer.Address)

	// RUN SUBSCRIBER
	stansub.RunSubscriber(sc, *config, db, csh)

	WaitForInterrupt(sc, csh)
}

// WaitForInterrupt ожидает прерывания работы сервиса
func WaitForInterrupt(sc stan.Conn, csh *cache.Cache) {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			csh.ClearCache()
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
