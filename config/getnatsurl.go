package config

import (
	"github.com/nats-io/nats.go"
	"os"
)

// GetNatsURL get Nats URL from environment variable
func GetNatsURL() string {
	// get natsURL from environment, or use defaultURL instead
	natsURL, ok := os.LookupEnv("NATS_URL")
	if !ok {
		natsURL = nats.DefaultURL
	}
	return natsURL
}
