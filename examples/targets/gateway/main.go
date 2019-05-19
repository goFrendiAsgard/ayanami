package main

import (
	"fmt"
	nats "github.com/nats-io/nats.go"
	"log"
	"net/http"
	"os"
	"strconv"
)

var routes []string

func init() {
	routes = []string{
		"/cowsay",
		"/figlet",
		"/",
	}
}

func main() {
	// get natsURL from environment, or use defaultURL instead
	natsURL, ok := os.LookupEnv("NATS_URL")
	if !ok {
		natsURL = nats.DefaultURL
	}
	// get port from environment, or use 8080 instead
	port, ok := os.LookupEnv("GATEWAY_PORT")
	if !ok {
		port = "8080"
	}
	// get maxMultipartFromLimit from environment, or use 20480 instead
	var multipartFormLimit, defaultMultipartFormLimit int64
	defaultMultipartFormLimit = 20480
	multipartFormLimit = defaultMultipartFormLimit
	strMultipartFormLimit, ok := os.LookupEnv("GATEWAY_MULTIPART_FORM_LIMIT")
	if ok {
		var err error
		multipartFormLimit, err = strconv.ParseInt(strMultipartFormLimit, 10, 64)
		if err != nil {
			multipartFormLimit = defaultMultipartFormLimit
		}
	}
	// create handlers and start server
	for _, route := range routes {
		handler := CreateHandler(natsURL, multipartFormLimit, route)
		http.HandleFunc(route, handler)
	}
	log.Printf("Listening on %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
