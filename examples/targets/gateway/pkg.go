package main

import (
	"mime/multipart"
)

// RequestHeaderPkg RequestHeaderPkg contains RequestData
type RequestHeaderPkg struct {
	ID   string              `json:"id"`
	Data map[string][]string `json:"servicedata"`
}

// RequestContentLengthPkg RequestContentLengthPkg contains RequestData
type RequestContentLengthPkg struct {
	ID   string `json:"id"`
	Data int64  `json:"servicedata"`
}

// RequestHostPkg RequestHostPkg contains Host and unique UUID
type RequestHostPkg struct {
	ID   string `json:"id"`
	Data string `json:"servicedata"`
}

// RequestFormPkg RequestFormPkg contains Form and unique UUID
type RequestFormPkg struct {
	ID   string              `json:"id"`
	Data map[string][]string `json:"servicedata"`
}

// RequestPostFormPkg RequestPostFormPkg contains PostForm and unique UUID
type RequestPostFormPkg struct {
	ID   string              `json:"id"`
	Data map[string][]string `json:"servicedata"`
}

// RequestMultipartFormPkg RequestMultipartFormPkg contains MultipartForm and unique UUID
type RequestMultipartFormPkg struct {
	ID   string          `json:"id"`
	Data *multipart.Form `json:"servicedata"`
}

// RequestMethodPkg RequestMethodPkg contains Method and unique UUID
type RequestMethodPkg struct {
	ID   string `json:"id"`
	Data string `json:"servicedata"`
}

// RequestRequestURIPkg RequestRequestURIPkg contains RequestURI and unique UUID
type RequestRequestURIPkg struct {
	ID   string `json:"id"`
	Data string `json:"servicedata"`
}

// RequestRemoteAddrPkg RequestRemoteAddrPkg contains RemoteAddr and unique UUID
type RequestRemoteAddrPkg struct {
	ID   string `json:"id"`
	Data string `json:"servicedata"`
}

// RequestJSONBodyPkg RequestJSONBodyPkg contains JSONBody and unique UUID
type RequestJSONBodyPkg struct {
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"servicedata"`
}

// ResponseCodePkg ResponseCodePkg contains Code and unique UUID to trigger response
type ResponseCodePkg struct {
	ID   string `json:"id"`
	Data int    `json:"servicedata"`
}

// ResponseContentPkg ResponseContentPkg contains Content and unique UUID to trigger response
type ResponseContentPkg struct {
	ID   string `json:"id"`
	Data string `json:"servicedata"`
}
