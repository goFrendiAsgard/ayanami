package service

import (
	"fmt"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
	"strings"
	"sync"
)

// FlowEvent event flow
type FlowEvent struct {
	InputEvent  string
	OutputEvent string
	VarName     string                        // read from inputEvent, put into var if value is not exists, publish into outputEvent
	UseValue    bool                          // if true, will use `Value` instead of `pkg.Data`
	Value       interface{}                   // value to override `pkg.Data` if `UseValue` is true
	UseFunction bool                          // if true, will use pass either `Value` or `pkg.Data` into `Function`, and publish the result
	Function    func(interface{}) interface{} // preprocessor, accept `Value` or `pkg.Data` before publish
}

// FlowService single flow config
type FlowService struct {
	FlowName string
	Input    IOList
	Output   IOList
	Flows    FlowEvents
}

// FlowEvents list of FlowEvent
type FlowEvents []FlowEvent

// GetInputEvents get unique input events
func (flows FlowEvents) GetInputEvents() []string {
	inputEvents := []string{}
	for _, flow := range flows {
		if flow.InputEvent != "" {
			inputEvents = AppendUniqueString(flow.InputEvent, inputEvents)
		}
	}
	return inputEvents
}

// GetVarFlowByInputEvent get unique vars by inputEvent
func (flows FlowEvents) GetVarFlowByInputEvent(inputEvent string) map[string]FlowEvent {
	varFlows := make(map[string]FlowEvent)
	for _, flow := range flows {
		if flow.InputEvent == inputEvent {
			varFlows[flow.VarName] = flow
		}
	}
	return varFlows
}

// GetVarNamesByInputEvent get unique vars by inputEvent
func (flows FlowEvents) GetVarNamesByInputEvent(inputEvent string) []string {
	varNames := []string{}
	for _, flow := range flows {
		if flow.InputEvent == inputEvent {
			varNames = AppendUniqueString(flow.VarName, varNames)
		}
	}
	return varNames
}

// GetOutputEventByVarNames get unique outputEvent by varNames
func (flows FlowEvents) GetOutputEventByVarNames(varName string) []string {
	outputEvents := []string{}
	for _, flow := range flows {
		if flow.VarName == varName {
			if flow.OutputEvent != "" {
				outputEvents = AppendUniqueString(flow.OutputEvent, outputEvents)
			}
		}
	}
	return outputEvents
}

// NewFlow create new service for flow
func NewFlow(flowName string, broker msgbroker.CommonBroker, inputs, outputs []string, flows FlowEvents) CommonService {
	// populate inputs
	var serviceInputs []IO
	for _, inputName := range inputs {
		eventName := fmt.Sprintf("flow.%s.in.%s", flowName, inputName)
		serviceInputs = append(serviceInputs, IO{VarName: inputName, EventName: eventName})
		for _, flow := range flows {
			if flow.VarName == inputName && flow.InputEvent != "" {
				eventName := flow.InputEvent
				serviceInputs = append(serviceInputs, IO{VarName: inputName, EventName: eventName})
			}
		}
	}
	// populate outputs
	var serviceOutputs []IO
	for _, outputName := range outputs {
		eventName := fmt.Sprintf("flow.%s.out.%s", flowName, outputName)
		serviceOutputs = append(serviceOutputs, IO{VarName: outputName, EventName: eventName})
		for _, flow := range flows {
			if flow.VarName == outputName && flow.OutputEvent != "" {
				eventName := flow.OutputEvent
				serviceOutputs = append(serviceOutputs, IO{VarName: outputName, EventName: eventName})
			}
		}
	}
	// get errorEventName
	errorEventName := fmt.Sprintf("flow.%s.err.message", flowName)
	// get flowWrappedFunction
	wrappedFunction := createFlowWrapper(flowName, broker, flows, inputs, outputs)
	return CommonService{
		Input:          serviceInputs,
		Output:         serviceOutputs,
		ErrorEventName: errorEventName,
		Function:       wrappedFunction,
	}
}

func createFlowWrapper(flowName string, broker msgbroker.CommonBroker, flows FlowEvents, inputVarNames, outputVarNames []string) WrappedFunction {
	return func(vars Dictionary) (Dictionary, error) {
		var lock sync.RWMutex // protecting vars
		ID, err := CreateID()
		if err != nil {
			return nil, err
		}
		rawInputEvents := getFlowRawInputEvents(flows, inputVarNames)
		var allConsumerDeclared sync.WaitGroup
		allConsumerDeclared.Add(len(rawInputEvents))
		outputCompleted := make(chan bool, 1)
		for rawInputEventIndex := range rawInputEvents {
			rawInputEvent := rawInputEvents[rawInputEventIndex]
			inputEvent := fmt.Sprintf("%s.%s", ID, rawInputEvent)
			log.Printf("[INFO: flow.%s] Consuming from %s", flowName, inputEvent)
			broker.Consume(inputEvent,
				func(pkg servicedata.Package) { // consume success
					log.Printf("[INFO: flow.%s] Getting message from %s: %#v", flowName, inputEvent, pkg.Data)
					for varName, varFlow := range flows.GetVarFlowByInputEvent(rawInputEvent) {
						allConsumerDeclared.Wait()
						lock.RLock()
						varExists := vars.Has(varName)
						lock.RUnlock()
						if varExists {
							log.Printf("[INFO: flow.%s] `%s` already defined, no need to set", flowName, varName)
							continue
						}
						publishedData := pkg.Data
						if varFlow.UseValue {
							publishedData = varFlow.Value
						}
						if varFlow.UseFunction && varFlow.Function != nil {
							publishedData = varFlow.Function(publishedData)
						}
						log.Printf("[INFO: flow.%s] Set `%s` into: `%#v`", flowName, varName, publishedData)
						lock.Lock()
						vars.Set(varName, publishedData)
						lock.Unlock()
						lock.RLock()
						publishFlowVar(flowName, broker, ID, flows, outputVarNames, varName, vars)
						lock.RUnlock()
					}
					lock.RLock()
					notifyIfOutputCompleted(vars, outputVarNames, outputCompleted)
					lock.RUnlock()
				},
				func(err error) { // consume error
					log.Printf("[ERROR: flow.%s] Error: %s", flowName, err)
					outputCompleted <- true
				},
			)
			allConsumerDeclared.Done()
		}
		// set default values
		lock.RLock()
		defaultVars := getFlowDefaultVars(flows, vars)
		lock.RUnlock()
		for varName, value := range defaultVars {
			lock.Lock()
			if !vars.Has(varName) {
				log.Printf("[INFO: flow.%s] Internally set `%s` into: `%#v`", flowName, varName, value)
				vars.Set(varName, value)
			} else {
				log.Printf("[INFO: flow.%s] `%s` already defined, no need to internal set", flowName, varName)
			}
			lock.Unlock()
		}
		for varName := range vars {
			lock.RLock()
			publishFlowVar(flowName, broker, ID, flows, outputVarNames, varName, vars)
			lock.RUnlock()
		}
		lock.RLock()
		notifyIfOutputCompleted(vars, outputVarNames, outputCompleted)
		lock.RUnlock()
		<-outputCompleted
		lock.RLock()
		outputs := createOutputs(outputVarNames, vars)
		log.Printf("[INFO: flow.%s] Internal flow `%s` ended. Outputs are: `%#v`", flowName, ID, outputs)
		return outputs, err
	}
}

func createOutputs(outputVarNames []string, vars Dictionary) Dictionary {
	outputs := make(Dictionary)
	for _, outputName := range outputVarNames {
		outputs.Set(outputName, vars.Get(outputName))
	}
	return outputs
}

func getFlowRawInputEvents(flows FlowEvents, inputVarNames []string) []string {
	candidates := flows.GetInputEvents()
	rawInputEvents := []string{}
	for _, candidate := range candidates {
		candidatePass := true
		varNames := flows.GetVarNamesByInputEvent(candidate)
		for _, varName := range varNames {
			if IsStringInArray(varName, inputVarNames) {
				candidatePass = false
				break
			}
		}
		if candidatePass {
			rawInputEvents = append(rawInputEvents, candidate)
		}
	}
	return rawInputEvents
}

func getFlowDefaultVars(flows FlowEvents, vars Dictionary) map[string]interface{} {
	defaultVars := make(map[string]interface{})
	// determine candidates from flows
	candidates := flows.GetVarNamesByInputEvent("")
	for _, candidate := range candidates {
		candidatePass := true
		var value interface{}
		for _, flow := range flows {
			if flow.VarName == candidate && flow.InputEvent == "" {
				value = flow.Value
			} else if isSubVarOf(flow.VarName, candidate) {
				candidatePass = false
				break
			}
		}
		if candidatePass {
			defaultVars[candidate] = value
		}
	}
	// add another candidates from predefined vars
	for varName, value := range vars {
		defaultVars[varName] = value
	}
	return defaultVars
}

func publishFlowVar(flowName string, broker msgbroker.CommonBroker, ID string, flows FlowEvents, outputVarNames []string, currentVarName string, vars Dictionary) {
	// if var is part of outputVarNames, ignore it. CommonBroker will do the job
	if IsStringInArray(currentVarName, outputVarNames) {
		return
	}
	// varNames contains currentVarName and all it's sub variable's names
	varNames := []string{currentVarName}
	for _, flow := range flows {
		varName := flow.VarName
		if varName != currentVarName && isSubVarOf(currentVarName, varName) {
			varNames = append(varNames, varName)
		}
	}
	// for every varNames, get it's outputEvent and publish
	for _, varName := range varNames {
		for _, rawOutputEvent := range flows.GetOutputEventByVarNames(varName) {
			varValue := vars.Get(varName)
			pkg := servicedata.Package{ID: ID, Data: varValue}
			outputEvent := fmt.Sprintf("%s.%s", ID, rawOutputEvent)
			log.Printf("[INFO: flow.%s] Publish into `%s`: `%#v`", flowName, outputEvent, pkg)
			broker.Publish(outputEvent, pkg)
		}
	}
}

// isSubVarOf determine whether subVarName is sub variable of varName or not
func isSubVarOf(varName, subVarName string) bool {
	return strings.Index(subVarName, varName+".") == 0
}

func notifyIfOutputCompleted(vars Dictionary, outputVarNames []string, outputCompleted chan bool) {
	if vars.HasAll(outputVarNames) {
		outputCompleted <- true
	}
}
