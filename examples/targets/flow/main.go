package main

import (
	nats "github.com/nats-io/nats.go"
	"os"
)

var configs Configs

func init() {
	configs = Configs{
		"flow": SingleConfig{
			Input: StringDictionary{
				"trig.request.get./.out.form": "form",
			},
			Output: StringDictionary{
				"code":    "trig.response.get./.in.code",
				"content": "trig.response.get./.in.content",
			},
			Function: WrappedEcho,
		},
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
