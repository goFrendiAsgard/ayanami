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
	var err error
	var pkg servicedata.Package
	// get json body
	decoder := json.NewDecoder(r.Body)
	var JSONBody map[string]interface{}
	decoder.Decode(&JSONBody)
	// publish header
	pkg = servicedata.Package{ID: ID, Data: r.Header}
	publishPkg(broker, ID, method, route, "header", pkg)
	// publish contentLength
	pkg = servicedata.Package{ID: ID, Data: r.ContentLength}
	publishPkg(broker, ID, method, route, "contentLength", pkg)
	// publish host
	pkg = servicedata.Package{ID: ID, Data: r.Host}
	publishPkg(broker, ID, method, route, "host", pkg)
	// publish form
	pkg = servicedata.Package{ID: ID, Data: r.Form}
	publishPkg(broker, ID, method, route, "form", pkg)
	// publish postForm
	pkg = servicedata.Package{ID: ID, Data: r.PostForm}
	publishPkg(broker, ID, method, route, "postForm", pkg)
	// publish multipartForm
	pkg = servicedata.Package{ID: ID, Data: r.MultipartForm}
	publishPkg(broker, ID, method, route, "multipartForm", pkg)
	// publish method
	pkg = servicedata.Package{ID: ID, Data: r.Method}
	publishPkg(broker, ID, method, route, "method", pkg)
	// publish requestURI
	pkg = servicedata.Package{ID: ID, Data: r.RequestURI}
	publishPkg(broker, ID, method, route, "requestURI", pkg)
	// publish remoteAddr
	pkg = servicedata.Package{ID: ID, Data: r.RemoteAddr}
	publishPkg(broker, ID, method, route, "remoteAddr", pkg)
	// publish jsonBody
	pkg = servicedata.Package{ID: ID, Data: JSONBody}
	publishPkg(broker, ID, method, route, "JSONBody", pkg)
	// return
	return err
}

func publishPkg(broker msgbroker.CommonBroker, ID string, method string, route string, varName string, pkg servicedata.Package) {
	eventName := fmt.Sprintf("%s.trig.request.%s.%s.out.%s", ID, method, route, varName)
	broker.Publish(eventName, pkg)
}

func consumeCode(broker msgbroker.CommonBroker, ID string, method string, route string, chListen chan bool, chData chan int, chErr chan error) {
	eventName := fmt.Sprintf("%s.trig.response.%s.%s.in.code", ID, method, route)
	log.Printf("[INFO] Prepare to consume `%s`", eventName)
	broker.Consume(eventName, func(pkg servicedata.Package) {
		chData <- pkg.Data.(int)
		chErr <- nil
		log.Printf("[INFO] Extract code from `%s`: `%#v`", eventName, pkg.Data)

	})
	log.Printf("[INFO] Start to consume from `%s`", eventName)
	chListen <- true
}

func consumeContent(broker msgbroker.CommonBroker, ID string, method string, route string, chListen chan bool, chData chan string, chErr chan error) {
	eventName := fmt.Sprintf("%s.trig.response.%s.%s.in.content", ID, method, route)
	log.Printf("[INFO] Prepare to consume `%s`", eventName)
	broker.Consume(eventName, func(pkg servicedata.Package) {
		chData <- pkg.Data.(string)
		chErr <- nil
		log.Printf("[INFO] Extract content from `%s`: `%#v`", eventName, pkg.Data)

	})
	log.Printf("[INFO] Start to consume from `%s`", eventName)
	chListen <- true
	log.Printf("[INFO] Start to consume from `%s`", eventName)
	chListen <- true
}
