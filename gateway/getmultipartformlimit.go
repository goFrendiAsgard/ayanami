package gateway

import (
	"os"
	"strconv"
)

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
