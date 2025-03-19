package requests

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"cloud.google.com/go/pubsub"

	client "github.com/H3ALY/T360_POST/Client"
	config "github.com/H3ALY/T360_POST/Config"
	pubSub "github.com/H3ALY/T360_POST/Publishers"
	search "github.com/H3ALY/T360_POST/Services/Search"
)

// HandleRequest processes incoming requests
func HandleRequest(cfg *config.Config, pubSubClient *pubsub.Client, w http.ResponseWriter, r *http.Request) {
	var searchBody search.SearchBody
	if err := json.NewDecoder(r.Body).Decode(&searchBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received VRM: %s, Date: %s", searchBody.Vrm, searchBody.ContraventionDate)

	// Create a channel to receive results
	resultChan := make(chan client.Result, 4)
	var wg sync.WaitGroup

	search.PerformSearch(cfg, searchBody, resultChan, &wg)

	wg.Wait()
	close(resultChan)

	results := search.CollectResults(resultChan)

	// Publish results to Pub/Sub
	topicID := cfg.Google.PubSubTopic
	for _, pubSubResult := range results {
		clientResult := client.Result{
			Reference: pubSubResult.Reference,
			Endpoint:  pubSubResult.Endpoint,
			Response:  pubSubResult.Response,
			Error:     pubSubResult.Error,
		}

		if err := pubSub.PublishToPubSub(pubSubClient, topicID, clientResult); err != nil {
			log.Printf("Failed to publish to Pub/Sub: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}
