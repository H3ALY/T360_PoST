package pubsubservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	config "github.com/H3ALY/T360_POST/Config"
	"google.golang.org/api/option"
)

func InitializePubSubClient(cfg *config.Config) (*pubsub.Client, error) {
	ctx := context.Background()
	var client *pubsub.Client
	var err error

	_, err = os.Stat(cfg.Google.ServiceAccountPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("service account key file does not exist at path: %s", cfg.Google.ServiceAccountPath)
	}

	file, err := os.Open(cfg.Google.ServiceAccountPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open service account key file: %v", err)
	}
	defer file.Close()

	var serviceAccount struct {
		ProjectID string `json:"project_id"`
	}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&serviceAccount); err != nil {
		return nil, fmt.Errorf("failed to parse service account key JSON: %v", err)
	}
	if serviceAccount.ProjectID == "" {
		return nil, fmt.Errorf("the service account key JSON is missing the 'project_id' field")
	}

	if !cfg.Google.UsingCloud {
		log.Println("Using Pub/Sub emulator")

		emulatorAddr := fmt.Sprintf("%s:%d", cfg.Emulator.Host, cfg.Emulator.Port)
		log.Printf("Connecting to emulator at %s", emulatorAddr)

		err := os.Setenv("PUBSUB_EMULATOR_HOST", emulatorAddr)
		if err != nil {
			log.Fatalf("Error setting emulator environment variable: %v", err)
		}

		client, err = pubsub.NewClient(ctx, cfg.Emulator.ProjectId)
		if err != nil {
			return nil, fmt.Errorf("failed to create Pub/Sub client: %v", err)
		}
	} else {

		log.Println("Using Google Cloud Pub/Sub")

		client, err = pubsub.NewClient(ctx, serviceAccount.ProjectID, option.WithCredentialsFile(cfg.Google.ServiceAccountPath))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create Pub/Sub client: %v", err)
	}

	return client, nil
}
