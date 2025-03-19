package searchservices

import (
	"log"
	"sync"
	"time"

	client "github.com/H3ALY/T360_POST/Client"
	config "github.com/H3ALY/T360_POST/Config"
	pubSub "github.com/H3ALY/T360_POST/Publishers"
)

type SearchBody struct {
	Vrm               string `json:"vrm"`
	ContraventionDate string `json:"contravention_date"`
}

// SearchAPIEndpoints holds the API endpoints for performing searches
type SearchAPIEndpoints struct {
	AcmeLease    string
	LeaseCompany string
	FleetCompany string
	HireCompany  string
}

// PerformSearch performs API calls to all endpoints concurrently and returns the results
func PerformSearch(cfg *config.Config, searchBody SearchBody, resultChan chan<- client.Result, wg *sync.WaitGroup) {
	endpoints := []string{
		cfg.Endpoints.TestSearch.AcmeLease,
		cfg.Endpoints.TestSearch.LeaseCompany,
		cfg.Endpoints.TestSearch.FleetCompany,
		cfg.Endpoints.TestSearch.HireCompany,
	}

	// Concurrently call each API endpoint
	for _, apiEndpoint := range endpoints {
		if apiEndpoint == "" {
			continue
		}
		wg.Add(1)
		go func(endpoint string) {
			defer wg.Done()
			// Call the API with the search parameters (VRM and Contravention Date)
			client.CallAPIWithTimeout(endpoint, searchBody.Vrm, searchBody.ContraventionDate, 2*time.Second, resultChan)
		}(apiEndpoint)
	}
}

// CollectResults collects the results from the channel and prepares them for publication
func CollectResults(resultChan <-chan client.Result) []pubSub.Result {
	var results []pubSub.Result
	for clientResult := range resultChan {
		if clientResult.Error == "" {
			results = append(results, pubSub.Result{
				Reference: clientResult.Reference,
				Endpoint:  clientResult.Endpoint,
				Response:  clientResult.Response,
				Error:     clientResult.Error,
			})
		} else {
			log.Printf("Skipping result due to error: %v", clientResult.Error)
		}
	}
	return results
}
