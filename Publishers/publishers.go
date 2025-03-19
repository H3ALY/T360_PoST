package publishers

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	client "github.com/H3ALY/T360_POST/Client"
)

type Result struct {
	Reference string `json:"reference"`
	Endpoint  string `json:"endpoint"`
	Response  string `json:"response"`
	Error     string `json:"error,omitempty"`
}

func PublishToPubSub(pubSubClient *pubsub.Client, topicID string, result client.Result) error {
	topic := pubSubClient.Topic(topicID)

	exists, err := topic.Exists(context.Background())

	if err != nil {
		return fmt.Errorf("error checking if topic exists: %v", err)
	}

	if !exists {
		return fmt.Errorf("topic %s does not exist", topicID)
	}

	// Create the message data to send to Pub/Sub
	data := []byte(result.Response)
	publishResult := topic.Publish(context.Background(), &pubsub.Message{
		Data: data,
		Attributes: map[string]string{
			"reference": result.Reference,
			"endpoint":  result.Endpoint,
		},
	})

	_, err = publishResult.Get(context.Background())
	if err != nil {
		log.Printf("Failed to publish to Pub/Sub: %v", err)
		return err
	}

	log.Printf("Successfully published message to Pub/Sub topic %s", topicID)
	return nil
}
