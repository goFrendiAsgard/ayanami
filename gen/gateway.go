package gen

import (
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
	"log"
	"path/filepath"
)

// GatewayConfig configuration to generate Gateway
type GatewayConfig struct {
	ServiceName string
	PackageName string
	Routes      []string
	*generator.IOHelper
	generator.StringHelper
}

// Validate validating config
func (config *GatewayConfig) Validate() bool {
	log.Printf("[INFO] Validating %s", config.ServiceName)
	if config.IsAlphaNumeric(config.ServiceName) {
		log.Printf("[ERROR] Service name should be alphanumeric, but `%s` found", config.ServiceName)
		return false
	}
	if config.PackageName == "" {
		log.Printf("[ERROR] Package name should not be empty")
		return false
	}
	return true
}

// Scaffold scaffolding config
func (config *GatewayConfig) Scaffold() error {
	return nil
}

// Build building config
func (config *GatewayConfig) Build() error {
	log.Printf("[INFO] Building %s", config.ServiceName)
	dirPath := fmt.Sprintf("gateway-%s", config.ServiceName)
	// write main.go
	log.Println("[INFO] Create main.go")
	mainPath := filepath.Join(dirPath, "main.go")
	err := config.WriteDep(mainPath, "gateway.main.go", config)
	if err != nil {
		return err
	}
	log.Println("[INFO] Create go.mod")
	goModPath := filepath.Join(dirPath, "go.mod")
	err = config.WriteDep(goModPath, "go.mod", config.PackageName)
	return err
}

// NewGateway create new gateway
func NewGateway(ioHelper *generator.IOHelper, serviceName string, packageName string, routes []string) GatewayConfig {
	return GatewayConfig{
		ServiceName: serviceName,
		PackageName: packageName,
		Routes:      routes,
		IOHelper:    ioHelper,
	}
}
