package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

type ServerConfig struct {
	Host     string
	Port     int
	DbID     int
	Password string
	Count    int
	WaitDur  int
	Queue    string
	Prefix   string
}

// Server A Snowstorm Server
type Server struct {
	redis  *redis.Client
	errors chan error
	count  int
	wait   time.Duration
	queue  string
	prefix string
}

// NewServer Create a new Snowstorm Server
func NewServer(config ServerConfig) *Server {
	server := new(Server)
	server.errors = make(chan error)
	server.redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DbID,
	})
	server.wait = time.Duration(config.WaitDur)
	server.count = config.Count
	server.queue = config.Queue
	server.prefix = config.Prefix
	if server.count == 0 {
		panic("Can't start server without at least 1 listener")
	}
	if server.count > 20 {
		panic("Can't start server with more than 20 listeners")
	}
	return server
}

// Run Begin waiting for requests.
//
//   count: The number of listeners to bootup.
func (server *Server) Run() {
	server.ensureRedisConnection()

	for i := 0; i < server.count; i++ {
		listener := NewListener(i, server.redis, server.queue, server.errors)
		go listener.Start(server.wait)
		go HandleMessages(server.prefix, server.redis, listener.messages, server.errors)
	}

	err := <-server.errors
	panic(err)
}

func (server *Server) ensureRedisConnection() {
	pong, err := server.redis.Ping().Result()
	fmt.Fprintf(os.Stderr, "checking redis: ")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stderr, "%s\n", pong)
}
