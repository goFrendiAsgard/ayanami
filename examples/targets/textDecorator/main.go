package main

import (
	nats "github.com/nats-io/nats.go"
	"os"
)

const serviceName = "textDecorator"

var configs Configs

func init() {
	configs = Configs{
		"pre": NewServiceConfig(
			serviceName,
			"pre",
			WrappedPre,
			[]string{"text"},
			[]string{"text"},
		),
		"cowsay": NewServiceConfig(
			serviceName,
			"cowsay",
			WrappedCowsay,
			[]string{"text"},
			[]string{"text"},
		),
		"figlet": NewServiceConfig(
			serviceName,
			"figlet",
			WrappedFiglet,
			[]string{"text"},
			[]string{"text"},
		),
	}
}

func main() {
	// get natsURL from environment, or use defaultURL instead
	natsURL, ok := os.LookupEnv("NATS_URL")
	if !ok {
		natsURL = nats.DefaultURL
	}
	// consume and publish forever
	ch := make(chan bool)
	ConsumeAndPublish(natsURL, configs)
	<-ch
}
