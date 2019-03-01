package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	versionMajor = 1
	versionMinor = 0
	versionPatch = 0
)

func main() {
	var config ServerConfig

	flag.StringVar(&config.Host, "host", "localhost", "redis host")
	flag.StringVar(&config.Queue, "queue", "snowstorm-generate", "redis key for inbound requests")
	flag.StringVar(&config.Prefix, "prefix", "snowstorm-delivery", "delivery queue prefix")
	flag.IntVar(&config.Port, "port", 6379, "redis port")
	flag.IntVar(&config.DbID, "db", 0, "redis db")
	flag.IntVar(&config.Count, "count", 1, "number of listeners")
	flag.IntVar(&config.WaitDur, "wait-ms", 5, "ms to wait before polling again.")

	flag.Parse()

	fmt.Fprintf(os.Stderr, "Snowstorm: %d.%d.%d\n", versionMajor, versionMinor, versionPatch)
	fmt.Fprintf(os.Stderr, "Redis: redis://%s:%d/%d\n", config.Host, config.Port, config.DbID)
	fmt.Fprintf(os.Stderr, "Queue: %s\n", config.Queue)
	fmt.Fprintf(os.Stderr, "Delivery: %s.{{ .ClientID }}.{{ .RequestID }}\n", config.Prefix)
	fmt.Fprintf(os.Stderr, "Listeners: %d\n", config.Count)
	fmt.Fprintf(os.Stderr, "Wait MS: %d\n", config.WaitDur)

	server := NewServer(config)
	server.Run()
}
