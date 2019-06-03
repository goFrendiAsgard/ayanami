package config

import (
	"os"
	"strconv"
)

// GetGatewayPort get port from environment
func GetGatewayPort() int64 {
	portStr, ok := os.LookupEnv("GATEWAY_PORT")
	if ok {
		port, err := strconv.ParseInt(portStr, 10, 64)
		if err != nil {
			return port
		}
	}
	return 8080
}
