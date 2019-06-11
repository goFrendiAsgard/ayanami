package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/service"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// Serve handle HTTP request
func Serve(broker msgbroker.CommonBroker, port int64, multipartFormLimit int64, routes []string) {
	sort.Sort(sort.Reverse(RouteSorter(routes)))
	log.Printf("[INFO: Gateway] Routes `%#v`", routes)
	for _, route := range routes {
		handler := createRouteHandler(broker, multipartFormLimit, route)
		http.HandleFunc(route, handler)
	}
	log.Printf("Listening on %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func createRouteHandler(broker msgbroker.CommonBroker, multipartFormLimit int64, route string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		method := strings.ToLower(r.Method)
		// create ID
		ID, err := service.CreateID()
		if err != nil {
			sendErrorResponse(ID, broker, method, route, w)
			return
		}
		// prepare channels
		responseChannel := make(chan map[string]interface{})
		consumeFromResponseTrigger(broker, ID, method, route, responseChannel)
		// publishToRequestTrigger
		err = publishToRequestTrigger(broker, ID, method, route, multipartFormLimit, r)
		if err != nil {
			sendErrorResponse(ID, broker, method, route, w)
			return
		}
		// wait for response
		response := <-responseChannel
		// response
		sendResponse(ID, broker, method, route, w, response)
	}
}

func createErrorResponse() map[string]interface{} {
	return map[string]interface{}{
		"code":    500,
		"content": "Internal Server Error",
		"header":  make(map[string]string),
	}
}

func consumeFromResponseTrigger(broker msgbroker.CommonBroker, ID, method, route string, responseChannel chan (map[string]interface{})) {
	eventName := getResponseEventName(ID, method, route)
	log.Printf("[INFO: Gateway] Subscribe `%s`", eventName)
	broker.Subscribe(eventName,
		// success
		func(pkg servicedata.Package) {
			log.Printf("[INFO: Gateway] Getting message from `%s`: `%#v`", eventName, pkg.Data)
			val, ok := pkg.Data.(map[string]interface{})
			if !ok {
				log.Printf("[ERROR: Gateway] Getting error while parsing data from `%s`", eventName)
				val = createErrorResponse()
			}
			responseChannel <- val
		},
		// error
		func(err error) {
			log.Printf("[ERROR: Gateway] Getting error while parsing `%s`: %s", eventName, err)
			responseChannel <- createErrorResponse()
		},
	)
}

func sendErrorResponse(ID string, broker msgbroker.CommonBroker, method, route string, w http.ResponseWriter) {
	response := createErrorResponse()
	sendResponse(ID, broker, method, route, w, response)
}

func sendResponse(ID string, broker msgbroker.CommonBroker, method, route string, w http.ResponseWriter, response map[string]interface{}) {
	eventName := getResponseEventName(ID, method, route)
	err := broker.Unsubscribe(eventName)
	if err != nil {
		log.Printf("[ERROR: Gateway] Getting error while unsubscribe from `%s`: %s", eventName, err)
		response = createErrorResponse()
	}
	// get code
	code := 200
	if codeInterface, exists := response["code"]; exists {
		var err error
		code, err = strconv.Atoi(fmt.Sprintf("%v", codeInterface))
		if err != nil {
			log.Printf("[ERROR: Gateway] Getting error while parsing code from `%s`: %s", eventName, err)
			code = 500
		}
	}
	// get content
	content := ""
	if code == 500 {
		content = "Internal Server Error"
	} else if contentInterface, exists := response["content"]; exists {
		content = fmt.Sprintf("%s", contentInterface)
	}
	// get header
	headers := map[string]string{
		"Content-Type": "text/html; charset=utf-8",
	}
	if headerInterface, exists := response["header"]; exists {
		var ok bool
		headers, ok = headerInterface.(map[string]string)
		if !ok {
			log.Printf("[ERROR: Gateway] Getting error while parsing header from `%s`: %#v", eventName, headerInterface)
		}
	}
	// set header, code, and content
	for key, val := range headers {
		w.Header().Set(key, val)
	}
	w.WriteHeader(code)
	_, err = fmt.Fprintf(w, "%s", content)
	if err != nil {
		log.Printf("[ERROR: Gateway] responding to %s: %s", ID, err)
	}
}

func publishToRequestTrigger(broker msgbroker.CommonBroker, ID string, method string, route string, multipartFormLimit int64, r *http.Request) error {
	eventName := getRequestEventName(ID, method, route)
	data := getDataForPublish(ID, multipartFormLimit, r)
	return service.Publish("Gateway", "", broker, ID, eventName, data)
}

func getDataForPublish(ID string, multipartFormLimit int64, r *http.Request) map[string]interface{} {
	// parse form & multipart form
	err := r.ParseForm()
	if err != nil {
		log.Printf("[ERROR: Gateway] Processing `%s`: %s", ID, err)
	}
	err = r.ParseMultipartForm(multipartFormLimit)
	if err != nil {
		log.Printf("[ERROR: Gateway] Processing `%s`: %s", ID, err)
	}
	data := make(map[string]interface{})
	data["method"] = r.Method
	data["URL"] = r.URL
	data["proto"] = r.Proto
	data["protoMajor"] = r.ProtoMajor
	data["protoMinor"] = r.ProtoMinor
	data["header"] = r.Header
	data["contentLength"] = r.ContentLength
	data["transferEncoding"] = r.TransferEncoding
	data["host"] = r.Host
	data["form"] = r.Form
	data["postForm"] = r.PostForm
	data["multipartForm"] = r.MultipartForm
	data["trailer"] = r.Trailer
	data["remoteAddr"] = r.RemoteAddr
	data["requestURI"] = r.RequestURI
	data["cookies"] = r.Cookies()
	data["userAgent"] = r.UserAgent()
	// get json body
	decoder := json.NewDecoder(r.Body)
	JSONBody := make(map[string]interface{})
	err = decoder.Decode(&JSONBody)
	if err != nil {
		log.Printf("[INFO: Gateway] Processing `%s`, request.body is not a valid JSON", ID)
	}
	data["JSONBody"] = JSONBody
	return data
}

func getResponseEventName(ID, method, route string) string {
	return getBaseEventName(ID, "response", method, route, "in")
}

func getRequestEventName(ID, method, route string) string {
	return getBaseEventName(ID, "request", method, route, "out")
}

func getBaseEventName(ID, trigger, method, route, direction string) string {
	segment := RouteToSegments(route)
	if segment == "" {
		return fmt.Sprintf("%s.trig.%s.%s.%s", ID, trigger, method, direction)
	}
	return fmt.Sprintf("%s.trig.%s.%s.%s.%s", ID, trigger, method, segment, direction)
}
