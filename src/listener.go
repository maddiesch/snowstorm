package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

type Listener struct {
	ID       int
	redis    *redis.Client
	messages chan Message
	errors   chan error
	key      string
}

func NewListener(ID int, redis *redis.Client, key string, errors chan error) *Listener {
	listener := new(Listener)
	listener.ID = ID
	listener.redis = redis
	listener.messages = make(chan Message)
	listener.errors = errors
	listener.key = key
	return listener
}

func (listener *Listener) Start(wait time.Duration) {
	fmt.Fprintf(os.Stderr, "Listener %d starting...\n", listener.ID)

	for {
		rawMessage, err := listener.redis.LPop(listener.key).Result()
		if err == redis.Nil {
			time.Sleep(wait * time.Millisecond)
		} else if err != nil {
			listener.errors <- err
		} else {
			processMessage(rawMessage, listener.messages)
		}
	}
}

func processMessage(raw string, messages chan Message) {
	var msg Message
	err := json.Unmarshal([]byte(raw), &msg)
	if err != nil {
		fmt.Println("Message failed to parse:", err)
	} else {
		messages <- msg
	}
}
