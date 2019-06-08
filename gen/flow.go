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
	generator.IOHelper
	generator.StringHelper
}

// Validate validating config
func (c FlowConfig) Validate() bool {
	serviceName := c.getServiceName()
	log.Printf("[INFO] Validating %s", serviceName)
	for _, input := range c.Inputs {
		if !c.IsMatch(input, "^[A-Za-z][a-zA-Z0-9]*$") {
			log.Printf("[ERROR] Invalid input `%s`", input)
			return false
		}
	}
	for _, output := range c.Outputs {
		if !c.IsMatch(output, "^[A-Za-z][a-zA-Z0-9]*$") {
			log.Printf("[ERROR] Invalid output `%s`", output)
			return false
		}
	}
	if !c.IsAlphaNumeric(c.FlowName) {
		log.Printf("[ERROR] Flow name should be alphanumeric, but `%s` found", c.FlowName)
		return false
	}
	if c.RepoName == "" {
		log.Printf("[ERROR] Repo name should not be empty")
		return false
	}
	for index, event := range c.Events {
		log.Printf("[INFO] Validating event %d", index)
		if !event.Validate() {
			return false
		}
	}
	return true
}

// Scaffold scaffolding config
func (c FlowConfig) Scaffold() error {
	serviceName := c.getServiceName()
	log.Printf("[INFO] Scaffolding %s", serviceName)
	for _, event := range c.Events {
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
		if !c.IsSourceExists(functionSourcePath) {
			log.Printf("[INFO] Create %s", functionFileName)
			err := c.WriteSource(functionSourcePath, "flow.function.go", data)
			if err != nil {
				return err
			}
		} else {
			log.Printf("[INFO] %s exists", functionFileName)
		}
		// write dependencies
		for _, dependency := range event.FunctionDependencies {
			dependencySourcePath := filepath.Join(packageSourcePath, dependency)
			if !c.IsSourceExists(dependencySourcePath) {
				log.Printf("[INFO] Create %s", dependency)
				err := c.WriteSource(dependencySourcePath, "dependency.go", data)
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
func (c FlowConfig) Build() error {
	serviceName := c.getServiceName()
	log.Printf("[INFO] Building %s", serviceName)
	depPath := fmt.Sprintf("flow-%s", c.FlowName)
	// write functions and dependencies
	for _, event := range c.Events {
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
		err := c.CopySourceToDep(functionSourcePath, functionDepPath)
		if err != nil {
			return err
		}
		// copy dependencies
		for _, dependency := range event.FunctionDependencies {
			log.Printf("[INFO] Create %s", dependency)
			dependencySourcePath := filepath.Join(packageSourcePath, dependency)
			dependencyDepPath := filepath.Join(packageDepPath, dependency)
			err := c.CopySourceToDep(dependencySourcePath, dependencyDepPath)
			if err != nil {
				return err
			}
		}
	}
	// write main.go
	log.Println("[INFO] Create main.go")
	mainPath := filepath.Join(depPath, "main.go")
	err := c.WriteDep(mainPath, "flow.main.go", c.toExposed())
	if err != nil {
		return err
	}
	// write go.mod
	log.Println("[INFO] Create go.mod")
	goModPath := filepath.Join(depPath, "go.mod")
	err = c.WriteDep(goModPath, "go.mod", c)
	return err
}

// CreateProgram create main.go and others
func (c FlowConfig) CreateProgram(depPath, serviceName, repoName, mainFunction string) {
	// TODO use this
}

// AddEvent add input to inputEvents
func (c *FlowConfig) AddEvent(event Event) {
	c.Events = append(c.Events, event)
}

// AddInputEvent create new Event
func (c *FlowConfig) AddInputEvent(eventName, varName string) {
	c.AddEvent(NewInputEvent(eventName, varName))
}

// AddOutputEvent create new Event
func (c *FlowConfig) AddOutputEvent(eventName, varName string) {
	c.AddEvent(NewOutputEvent(eventName, varName))
}

// AddOutputEventVal create new Event with value
func (c *FlowConfig) AddOutputEventVal(eventName, varName string, value interface{}) {
	c.AddEvent(NewOutputEventVal(eventName, varName, value))
}

// AddOutputEventFunc create new Event with function
func (c *FlowConfig) AddOutputEventFunc(eventName, varName, functionPackage, functionName string, functionDependencies []string) {
	c.AddEvent(NewOutputEventFunc(eventName, varName, functionPackage, functionName, functionDependencies))
}

func (c *FlowConfig) toExposed() ExposedFlowConfig {
	return ExposedFlowConfig{
		ServiceName: c.getServiceName(),
		RepoName:    c.RepoName,
		FlowName:    c.FlowName,
		Packages:    c.getPackagesForExposed(),
		Events:      c.getEventsForExposed(),
		Outputs:     c.QuoteArrayAndJoin(c.Outputs, ", "),
		Inputs:      c.QuoteArrayAndJoin(c.Inputs, ", "),
	}
}

func (c *FlowConfig) getServiceName() string {
	return fmt.Sprintf("flow%s", c.FlowName)
}

func (c *FlowConfig) getEventsForExposed() []map[string]string {
	events := []map[string]string{}
	for _, event := range c.Events {
		events = append(events, event.ToMap())
	}
	return events
}

func (c *FlowConfig) getPackagesForExposed() []string {
	packages := []string{}
	for _, event := range c.Events {
		if !event.UseFunction {
			continue
		}
		packages = append(packages, event.FunctionPackage)
	}
	return packages
}

// NewFlow create new flow
func NewFlow(g *generator.Generator, repoName, flowName string, inputs, outputs []string, events []Event) FlowConfig {
	return FlowConfig{
		RepoName: repoName,
		FlowName: flowName,
		Inputs:   inputs,
		Outputs:  outputs,
		Events:   events,
		IOHelper: g.IOHelper,
	}
}

// NewEmptyFlow create new empty flow
func NewEmptyFlow(g *generator.Generator, repoName, flowName string, inputs, outputs []string) FlowConfig {
	return NewFlow(g, repoName, flowName, inputs, outputs, []Event{})
}
