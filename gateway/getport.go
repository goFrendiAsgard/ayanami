package gateway

import (
	"os"
	"strconv"
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
