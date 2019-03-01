package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/segmentio/ksuid"
)

// Message A message received from the client.
type Message struct {
	// ClientID The ID of the requestor
	ClientID string `json:"ClientID"`

	// RequestID The ID for each generate request
	RequestID string `json:"RequestID"`
}

type response struct {
	GeneratedID string `json:"GeneratedID"`
	RequestID   string `json:"RequestID"`
}

// HandleMessages Handles messages as they come in from the listener.
func HandleMessages(prefix string, client *redis.Client, messages chan Message, errors chan error) {
	for {
		msg := <-messages
		kid, err := ksuid.NewRandom()
		if err != nil {
			errors <- err
		} else {
			id := kid.String()
			go msg.deliver(client, prefix, id, errors)
		}
	}
}

func (msg *Message) deliver(client *redis.Client, prefix string, id string, errors chan error) {
	fmt.Fprintf(os.Stderr, "deliver: %s %s %s\n", id, msg.ClientID, msg.RequestID)
	key := fmt.Sprintf("%s.%s.%s", prefix, msg.ClientID, msg.RequestID)
	result := &response{id, msg.RequestID}
	payload, err := json.Marshal(result)
	if err != nil {
		errors <- err
	} else {
		client.LPush(key, payload)
		client.Expire(key, 10*time.Second)
	}
}
