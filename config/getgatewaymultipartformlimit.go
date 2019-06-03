package config

import (
	"os"
	"strconv"
)

// GetGatewayMultipartFormLimit get port from environment
func GetGatewayMultipartFormLimit() int64 {
	multipartFormLimitStr, ok := os.LookupEnv("GATEWAY_MULTIPART_FORM_LIMIT")
	if ok {
		multipartFormLimit, err := strconv.ParseInt(multipartFormLimitStr, 10, 64)
		if err != nil {
			return multipartFormLimit
		}
	}
	return 20480
}
