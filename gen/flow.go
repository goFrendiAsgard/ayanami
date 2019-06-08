package gen

import (
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
	"log"
	"path/filepath"
)

// ExposedFlowConfig exposed ready flowConfig
type ExposedFlowConfig struct {
	ServiceName string
	RepoName    string
	FlowName    string
	Packages    []string
	Events      []map[string]string
	Outputs     string
	Inputs      string
}

// FlowConfig definition
type FlowConfig struct {
	RepoName string
	FlowName string
	Inputs   []string
	Outputs  []string
	Events   []Event
	*generator.IOHelper
	generator.StringHelper
}

// AddEvent add input to inputEvents
func (config *FlowConfig) AddEvent(event Event) {
	config.Events = append(config.Events, event)
}

// Validate validating config
func (config FlowConfig) Validate() bool {
	log.Printf("[INFO] Validating %s", config.FlowName)
	for _, input := range config.Inputs {
		if !config.IsMatch(input, "^[A-Za-z][a-zA-Z0-9]*$") {
			log.Printf("[ERROR] Invalid input `%s`", input)
			return false
		}
	}
	for _, output := range config.Outputs {
		if !config.IsMatch(output, "^[A-Za-z][a-zA-Z0-9]*$") {
			log.Printf("[ERROR] Invalid output `%s`", output)
			return false
		}
	}
	if !config.IsAlphaNumeric(config.FlowName) {
		log.Printf("[ERROR] Flow name should be alphanumeric, but `%s` found", config.FlowName)
		return false
	}
	if config.RepoName == "" {
		log.Printf("[ERROR] Repo name should not be empty")
		return false
	}
	for index, event := range config.Events {
		log.Printf("[INFO] Validating event %d", index)
		if !event.Validate() {
			return false
		}
	}
	return true
}

// Scaffold scaffolding config
func (config FlowConfig) Scaffold() error {
	log.Printf("[INFO] Scaffolding %s", config.FlowName)
	for _, event := range config.Events {
		if !event.UseFunction {
			continue
		}
		data := map[string]string{
			"FunctionPackage": event.FunctionPackage,
			"FunctionName":    event.FunctionName,
		}
		packageSourcePath := event.FunctionPackage
		functionFileName := event.GetFunctionFileName()
		// write function
		functionSourcePath := filepath.Join(packageSourcePath, functionFileName)
		if !config.IsSourceExists(functionSourcePath) {
			log.Printf("[INFO] Create %s", functionFileName)
			err := config.WriteSource(functionSourcePath, "flow.function.go", data)
			if err != nil {
				return err
			}
		} else {
			log.Printf("[INFO] %s exists", functionFileName)
		}
		// write dependencies
		for _, dependency := range event.FunctionDependencies {
			dependencySourcePath := filepath.Join(packageSourcePath, dependency)
			if !config.IsSourceExists(dependencySourcePath) {
				log.Printf("[INFO] Create %s", dependency)
				err := config.WriteSource(dependencySourcePath, "dependency.go", data)
				if err != nil {
					return err
				}
			} else {
				log.Printf("[INFO] %s exists", dependency)
			}
		}
	}
	return nil
}

// Build building config
func (config FlowConfig) Build() error {
	log.Printf("[INFO] Building %s", config.FlowName)
	depPath := fmt.Sprintf("flow-%s", config.FlowName)
	// write functions and dependencies
	for _, event := range config.Events {
		if !event.UseFunction {
			continue
		}
		packageSourcePath := event.FunctionPackage
		packageDepPath := filepath.Join(depPath, event.FunctionPackage)
		functionFileName := event.GetFunctionFileName()
		// copy function
		log.Printf("[INFO] Create %s", functionFileName)
		functionSourcePath := filepath.Join(packageSourcePath, functionFileName)
		functionDepPath := filepath.Join(packageDepPath, functionFileName)
		err := config.CopySourceToDep(functionSourcePath, functionDepPath)
		if err != nil {
			return err
		}
		// copy dependencies
		for _, dependency := range event.FunctionDependencies {
			log.Printf("[INFO] Create %s", dependency)
			dependencySourcePath := filepath.Join(packageSourcePath, dependency)
			dependencyDepPath := filepath.Join(packageDepPath, dependency)
			err := config.CopySourceToDep(dependencySourcePath, dependencyDepPath)
			if err != nil {
				return err
			}
		}
	}
	// write main.go
	log.Println("[INFO] Create main.go")
	mainPath := filepath.Join(depPath, "main.go")
	err := config.WriteDep(mainPath, "flow.main.go", config.toExposed())
	if err != nil {
		return err
	}
	// write go.mod
	log.Println("[INFO] Create go.mod")
	goModPath := filepath.Join(depPath, "go.mod")
	err = config.WriteDep(goModPath, "go.mod", config)
	return err
}

func (config *FlowConfig) toExposed() ExposedFlowConfig {
	return ExposedFlowConfig{
		ServiceName: fmt.Sprintf("flow%s", config.FlowName),
		RepoName:    config.RepoName,
		FlowName:    config.FlowName,
		Packages:    config.getPackagesForExposed(),
		Events:      config.getEventsForExposed(),
		Outputs:     config.QuoteArrayAndJoin(config.Outputs, ", "),
		Inputs:      config.QuoteArrayAndJoin(config.Inputs, ", "),
	}
}

func (config *FlowConfig) getEventsForExposed() []map[string]string {
	events := []map[string]string{}
	for _, event := range config.Events {
		events = append(events, event.ToMap())
	}
	return events
}

func (config *FlowConfig) getPackagesForExposed() []string {
	packages := []string{}
	for _, event := range config.Events {
		if !event.UseFunction {
			continue
		}
		packages = append(packages, event.FunctionPackage)
	}
	return packages
}

// NewFlow create new flow
func NewFlow(ioHelper *generator.IOHelper, repoName, flowName string, inputs, outputs []string, events []Event) FlowConfig {
	return FlowConfig{
		RepoName: repoName,
		FlowName: flowName,
		Inputs:   inputs,
		Outputs:  outputs,
		Events:   events,
		IOHelper: ioHelper,
	}
}

// NewEmptyFlow create new empty flow
func NewEmptyFlow(ioHelper *generator.IOHelper, repoName, flowName string, inputs, outputs []string) FlowConfig {
	return NewFlow(ioHelper, repoName, flowName, inputs, outputs, []Event{})
}
