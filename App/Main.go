package main

import (
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/pubsub"

	config "github.com/H3ALY/T360_POST/Config"
	requests "github.com/H3ALY/T360_POST/Models"
	pubSubService "github.com/H3ALY/T360_POST/Services/PubSub"
)

var pubSubClient *pubsub.Client

func main() {
	// Load the configuration from config.yaml
	cfg, err := config.LoadConfig("Config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	pubSubClient, err = pubSubService.InitializePubSubClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}
	defer pubSubClient.Close()

	// Read the port from the config
	port := fmt.Sprintf(":%d", cfg.Server.Port)

	// Define the route and associate it with the handler function
	http.HandleFunc("/handle_request", func(w http.ResponseWriter, r *http.Request) {
		requests.HandleRequest(cfg, pubSubClient, w, r)
	})

	// Start the HTTP server with the configured port
	log.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
