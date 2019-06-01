package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/service"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// GetPort get port from environment
func GetPort() int64 {
	portStr, ok := os.LookupEnv("GATEWAY_PORT")
	if ok {
		port, err := strconv.ParseInt(portStr, 10, 64)
		if err != nil {
			return port
		}
	}
	return 8080
}

// GetMultipartFormLimit get port from environment
func GetMultipartFormLimit() int64 {
	multipartFormLimitStr, ok := os.LookupEnv("GATEWAY_MULTIPART_FORM_LIMIT")
	if ok {
		multipartFormLimit, err := strconv.ParseInt(multipartFormLimitStr, 10, 64)
		if err != nil {
			return multipartFormLimit
		}
	}
	return 20480
}

// Serve handle HTTP request
func Serve(broker msgbroker.CommonBroker, port int64, multipartFormLimit int64, routes []string) {
	for _, route := range routes {
		handler := createRouteHandler(broker, multipartFormLimit, route)
		http.HandleFunc(route, handler)
	}
	log.Printf("Listening on %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func createRouteHandler(broker msgbroker.CommonBroker, multipartFormLimit int64, route string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// create ID
		ID, err := service.CreateID()
		if err != nil {
			responseError(ID, w, 500, err)
			return
		}
		method := strings.ToLower(r.Method)
		// prepare channels
		codeChannel := make(chan int, 1)
		contentChannel := make(chan string, 1)
		consume(broker, ID, method, route, codeChannel, contentChannel)
		// publish
		err = publish(broker, ID, method, route, multipartFormLimit, r)
		if err != nil {
			responseError(ID, w, 500, err)
			return
		}
		// wait for response
		code := <-codeChannel
		content := <-contentChannel
		response(ID, w, code, content)
	}
}

func consume(broker msgbroker.CommonBroker, ID, method, route string, codeChannel chan int, contentChannel chan string) {
	codeEventName := fmt.Sprintf("%s.trig.response.%s %s.in.code", ID, method, route)
	contentEventName := fmt.Sprintf("%s.trig.response.%s %s.in.content", ID, method, route)
	// consume code
	log.Printf("[INFO: Gateway] Consume `%s`", codeEventName)
	broker.Consume(codeEventName,
		// success
		func(pkg servicedata.Package) {
			log.Printf("[INFO: Gateway] Getting message from `%s`: `%#v`", codeEventName, pkg.Data)
			val, err := strconv.Atoi(fmt.Sprintf("%v", pkg.Data))
			if err == nil {
				codeChannel <- val
			} else {
				log.Printf("[ERROR: Gateway] %s", err)
			}
		},
		// error
		func(err error) {
			codeChannel <- 500
			contentChannel <- "Internal Server Error"
		},
	)
	// codeChannel <- 200
	// consume event
	log.Printf("[INFO: Gateway] Consume `%s`", contentEventName)
	broker.Consume(contentEventName,
		// success
		func(pkg servicedata.Package) {
			log.Printf("[INFO: Gateway] Getting message from `%s`: `%#v`", contentEventName, pkg.Data)
			contentChannel <- fmt.Sprintf("%s", pkg.Data)
		},
		// error
		func(err error) {
			codeChannel <- 500
			contentChannel <- "Internal Server Error"
		},
	)
}

func responseError(ID string, w http.ResponseWriter, code int, err error) {
	content := fmt.Sprintf("%s", err)
	response(ID, w, code, content)
}

func response(ID string, w http.ResponseWriter, code int, content string) {
	log.Printf("[INFO: Gateway] responding to %s: %d, %s", ID, code, content)
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", content)
}

func publish(broker msgbroker.CommonBroker, ID string, method string, route string, multipartFormLimit int64, r *http.Request) error {
	eventName := fmt.Sprintf("%s.trig.request.%s %s.out.req", ID, method, route)
	// parse form & multipart form
	r.ParseForm()
	r.ParseMultipartForm(multipartFormLimit)
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
		log.Printf("[INFO: Gateway] Processing `%s`, request.body is not a valid JSON", ID)
	}
	data["JSONBody"] = JSONBody
	// prepare pkg
	pkg := servicedata.Package{ID: ID, Data: data}
	log.Printf("[INFO: Gateway] publish `%s`: %#v", eventName, pkg)
	return broker.Publish(eventName, pkg)
}
