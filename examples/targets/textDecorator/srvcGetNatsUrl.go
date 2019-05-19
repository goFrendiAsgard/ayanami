package main

import (
	nats "github.com/nats-io/nats.go"
	"os"
)

// GetNatsURL getting nats URL from environment
func GetNatsURL() string {
	// get natsURL from environment, or use defaultURL instead
	natsURL, ok := os.LookupEnv("NATS_URL")
	if !ok {
		natsURL = nats.DefaultURL
	}
	return natsURL
}
