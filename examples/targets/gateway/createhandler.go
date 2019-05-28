package main

import (
	"encoding/json"
	"fmt"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/service"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
	"net/http"
	"strings"
)

// CreateHandler create handler function
func CreateHandler(multipartFormLimit int64, route string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// create ID
		ID, err := service.CreateID()
		if err != nil {
			responseError(ID, w, 500, err)
			return
		}
		// connect to nats
		broker, err := msgbroker.NewNats()
		if err != nil {
			responseError(ID, w, 500, err)
			return
		}
		// parse form & multipart form
		r.ParseForm()
		r.ParseMultipartForm(multipartFormLimit)
		method := strings.ToLower(r.Method)
		// listen to code
		chListenCode := make(chan bool)
		chListenCodeErr := make(chan error)
		chCode := make(chan int)
		go consumeCode(broker, ID, method, route, chListenCode, chCode, chListenCodeErr)
		// listen to content
		chListenContent := make(chan bool)
		chListenContentErr := make(chan error)
		chContent := make(chan string)
		go consumeContent(broker, ID, method, route, chListenContent, chContent, chListenContentErr)
		// start publish
		<-chListenCode
		<-chListenContent
		// publish
		err = publish(broker, ID, method, route, r)
		if err != nil {
			responseError(ID, w, 500, err)
			return
		}
		// get the code
		code := <-chCode
		err = <-chListenCodeErr
		if err != nil {
			responseError(ID, w, 500, err)
			return
		}
		// get the content
		content := <-chContent
		err = <-chListenContentErr
		if err != nil {
			responseError(ID, w, 500, err)
			return
		}
		// return normal response
		log.Printf("[INFO] Responding to %s: %d, %s", ID, code, content)
		w.WriteHeader(code)
		fmt.Fprintf(w, "%s", content)
	}
}

func responseError(ID string, w http.ResponseWriter, code int, err error) {
	log.Printf("[ERROR] responding to %s: %d, %s", ID, code, err)
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", err)
}

func publish(broker msgbroker.CommonBroker, ID string, method string, route string, r *http.Request) error {
	eventName := fmt.Sprintf("%s.trig.request.%s %s.out.req", ID, method, route)
	log.Printf("[INFO] prepare to publish `%s`", eventName)
	data := make(map[string]interface{})
	data["header"] = r.Header
	data["contentLength"] = r.ContentLength
	data["host"] = r.Host
	data["form"] = r.Form
	data["postForm"] = r.PostForm
	data["multipartForm"] = r.MultipartForm
	data["method"] = r.Method
	data["requestURI"] = r.RequestURI
	data["remoteAddr"] = r.RemoteAddr
	// get json body
	decoder := json.NewDecoder(r.Body)
	JSONBody := make(map[string]interface{})
	err := decoder.Decode(&JSONBody)
	if err != nil {
		log.Printf("[INFO] Fail to decode JSONBody while processing `%s`", eventName)
	}
	data["JSONBody"] = JSONBody
	// prepare pkg
	pkg := servicedata.Package{ID: ID, Data: data}
	log.Printf("[INFO] publish `%s`: %#v", eventName, pkg)
	return broker.Publish(eventName, pkg)
}

func consumeCode(broker msgbroker.CommonBroker, ID string, method string, route string, chListen chan bool, chData chan int, chErr chan error) {
	eventName := fmt.Sprintf("%s.trig.response.%s %s.in.code", ID, method, route)
	log.Printf("[INFO] Prepare to consume `%s`", eventName)
	broker.Consume(eventName,
		// success
		func(pkg servicedata.Package) {
			chData <- pkg.Data.(int)
			chErr <- nil
			log.Printf("[INFO] Extract code from `%s`: `%#v`", eventName, pkg.Data)
		},
		// error
		func(err error) {
			chData <- 0
			chErr <- nil
		},
	)
	log.Printf("[INFO] Start to consume from `%s`", eventName)
	chListen <- true
}

func consumeContent(broker msgbroker.CommonBroker, ID string, method string, route string, chListen chan bool, chData chan string, chErr chan error) {
	eventName := fmt.Sprintf("%s.trig.response.%s %s.in.content", ID, method, route)
	log.Printf("[INFO] Prepare to consume `%s`", eventName)
	broker.Consume(eventName,
		// success
		func(pkg servicedata.Package) {
			chData <- pkg.Data.(string)
			chErr <- nil
			log.Printf("[INFO] Extract content from `%s`: `%#v`", eventName, pkg.Data)
		},
		// error
		func(err error) {
			chData <- ""
			chErr <- nil
		},
	)
	log.Printf("[INFO] Start to consume from `%s`", eventName)
	chListen <- true
}
