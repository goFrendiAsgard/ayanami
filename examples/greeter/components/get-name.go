package main

import (
	"net/http"
)

func GetParam(request http.Request, paramName string) string {
	return ""
}

func GetName(request http.Request) string {
	return GetParam(request, "name")
}

func GetYear(request http.Request) string {
	return GetParam(request, "year")
}

func ComposeHttpResponse(string) http.Response {
	return http.Response{}
}
