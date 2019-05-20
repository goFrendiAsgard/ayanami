package main

import (
	"encoding/json"
	"fmt"
	nats "github.com/nats-io/nats.go"
	"log"
	"net/http"
	"strings"
)

// CreateHandler create handler function
func CreateHandler(multipartFormLimit int64, route string) func(http.ResponseWriter, *http.Request) {
	natsURL := SrvcGetNatsURL()
	return func(w http.ResponseWriter, r *http.Request) {
		// create ID
		ID, err := SrvcCreateID()
		if err != nil {
			responseError(ID, w, 500, err)
			return
		}
		// connect to nats
		nc, err := nats.Connect(natsURL)
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
		go consumeCode(nc, ID, method, route, chListenCode, chCode, chListenCodeErr)
		// listen to content
		chListenContent := make(chan bool)
		chListenContentErr := make(chan error)
		chContent := make(chan string)
		go consumeContent(nc, ID, method, route, chListenContent, chContent, chListenContentErr)
		// start publish
		<-chListenCode
		<-chListenContent
		// publish
		err = publish(nc, ID, method, route, r)
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
		// in gateway case, nats connection will only used until we get reply from flow
		nc.Drain()
		nc.Close()
	}
}

func responseError(ID string, w http.ResponseWriter, code int, err error) {
	log.Printf("[ERROR] responding to %s: %d, %s", ID, code, err)
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", err)
}

func publish(nc *nats.Conn, ID string, method string, route string, r *http.Request) error {
	var err error
	var pkg interface{}
	// get json body
	decoder := json.NewDecoder(r.Body)
	var JSONBody map[string]interface{}
	decoder.Decode(&JSONBody)
	// publish header
	pkg = RequestHeaderPkg{ID: ID, Data: r.Header}
	err = publishPkg(nc, ID, method, route, "header", pkg)
	if err != nil {
		return err
	}
	// publish contentLength
	pkg = RequestContentLengthPkg{ID: ID, Data: r.ContentLength}
	err = publishPkg(nc, ID, method, route, "contentLength", pkg)
	if err != nil {
		return err
	}
	// publish host
	pkg = RequestHostPkg{ID: ID, Data: r.Host}
	err = publishPkg(nc, ID, method, route, "host", pkg)
	if err != nil {
		return err
	}
	// publish form
	pkg = RequestFormPkg{ID: ID, Data: r.Form}
	err = publishPkg(nc, ID, method, route, "form", pkg)
	if err != nil {
		return err
	}
	// publish postForm
	pkg = RequestPostFormPkg{ID: ID, Data: r.PostForm}
	err = publishPkg(nc, ID, method, route, "postForm", pkg)
	if err != nil {
		return err
	}
	// publish multipartForm
	pkg = RequestMultipartFormPkg{ID: ID, Data: r.MultipartForm}
	err = publishPkg(nc, ID, method, route, "multipartForm", pkg)
	if err != nil {
		return err
	}
	// publish method
	pkg = RequestMethodPkg{ID: ID, Data: r.Method}
	err = publishPkg(nc, ID, method, route, "method", pkg)
	if err != nil {
		return err
	}
	// publish requestURI
	pkg = RequestRequestURIPkg{ID: ID, Data: r.RequestURI}
	err = publishPkg(nc, ID, method, route, "requestURI", pkg)
	if err != nil {
		return err
	}
	// publish remoteAddr
	pkg = RequestRequestURIPkg{ID: ID, Data: r.RemoteAddr}
	err = publishPkg(nc, ID, method, route, "remoteAddr", pkg)
	if err != nil {
		return err
	}
	// publish jsonBody
	pkg = RequestJSONBodyPkg{ID: ID, Data: JSONBody}
	err = publishPkg(nc, ID, method, route, "JSONBody", pkg)
	if err != nil {
		return err
	}
	// return
	return err
}

func publishPkg(nc *nats.Conn, ID string, method string, route string, varName string, pkg interface{}) error {
	eventName := fmt.Sprintf("%s.trig.request.%s.%s.out.%s", ID, method, route, varName)
	JSONByte, err := json.Marshal(&pkg)
	if err != nil {
		log.Printf("[ERROR] %s: %s", eventName, err)
		return err
	}
	log.Printf("[INFO] Publish into `%s`: `%#v`", eventName, pkg)
	nc.Publish(eventName, JSONByte)
	return nil
}

func consumeCode(nc *nats.Conn, ID string, method string, route string, chListen chan bool, chData chan int, chErr chan error) {
	eventName := fmt.Sprintf("%s.trig.response.%s.%s.in.code", ID, method, route)
	nc.Subscribe(eventName, func(m *nats.Msg) {
		log.Printf("[INFO] Get code from `%s`: `%s`", eventName, string(m.Data))
		var pkg ResponseCodePkg
		JSONByte := m.Data
		err := json.Unmarshal(JSONByte, &pkg)
		if err != nil {
			chData <- 0
			chErr <- err
			return
		}
		chData <- pkg.Data
		chErr <- nil
		log.Printf("[INFO] Extract code from `%s`: `%#v`", eventName, pkg.Data)
	})
	log.Printf("[INFO] Start to consume from `%s`", eventName)
	chListen <- true
}

func consumeContent(nc *nats.Conn, ID string, method string, route string, chListen chan bool, chData chan string, chErr chan error) {
	eventName := fmt.Sprintf("%s.trig.response.%s.%s.in.content", ID, method, route)
	log.Printf("[INFO] Prepare to consume `%s`", eventName)
	nc.Subscribe(eventName, func(m *nats.Msg) {
		log.Printf("[INFO] Get content from `%s`: `%s`", eventName, string(m.Data))
		var pkg ResponseContentPkg
		JSONByte := m.Data
		err := json.Unmarshal(JSONByte, &pkg)
		if err != nil {
			chData <- ""
			chErr <- err
			return
		}
		chData <- pkg.Data
		chErr <- nil
		log.Printf("[INFO] Extract content from `%s`: `%#v`", eventName, pkg.Data)
	})
	log.Printf("[INFO] Start to consume from `%s`", eventName)
	chListen <- true
}
