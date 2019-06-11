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
	VarName     string                        // read from inputEvent, put into var if value is not exists, publishServiceOutput into outputEvent
	UseValue    bool                          // if true, will use `Value` instead of `pkg.Data`
	Value       interface{}                   // value to override `pkg.Data` if `UseValue` is true
	UseFunction bool                          // if true, will use pass either `Value` or `pkg.Data` into `Function`, and publishServiceOutput the result
	Function    func(interface{}) interface{} // preprocessor, accept `Value` or `pkg.Data` before publishServiceOutput
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
	var inputEvents []string
	for _, flow := range flows {
		if flow.InputEvent != "" {
			inputEvents = AppendUniqueString(flow.InputEvent, inputEvents)
		}
	}
	return inputEvents
}

// GetInputEventByVarName get unique input events by var
func (flows FlowEvents) GetInputEventByVarName(varName string) []string {
	var inputEvents []string
	for _, flow := range flows {
		if flow.InputEvent != "" && flow.VarName == varName {
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
	var varNames []string
	for _, flow := range flows {
		if flow.InputEvent == inputEvent {
			varNames = AppendUniqueString(flow.VarName, varNames)
		}
	}
	return varNames
}

// GetOutputEventByVarNames get unique outputEvent by varNames
func (flows FlowEvents) GetOutputEventByVarNames(varName string) []string {
	var outputEvents []string
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
		ServiceName:    "flow",
		MethodName:     flowName,
		Input:          serviceInputs,
		Output:         serviceOutputs,
		ErrorEventName: errorEventName,
		Function:       wrappedFunction,
	}
}

func createFlowWrapper(flowName string, broker msgbroker.CommonBroker, flows FlowEvents, inputVarNames, outputVarNames []string) WrappedFunction {
	return func(presetVars Dictionary) (Dictionary, error) {
		var lock sync.RWMutex // protecting vars
		vars := make(Dictionary)
		ID, err := CreateID()
		if err != nil {
			return nil, err
		}
		rawInputEvents := getFlowRawInputEvents(flows, inputVarNames)
		var allConsumerDeclared sync.WaitGroup
		allConsumerDeclared.Add(len(rawInputEvents))
		outputCompleted := make(chan bool, 1)
		successHandlers := map[string]func(servicedata.Package){}
		for rawInputEventIndex := range rawInputEvents {
			rawInputEvent := rawInputEvents[rawInputEventIndex]
			inputEvent := fmt.Sprintf("%s.%s", ID, rawInputEvent)
			// create handlers & consume
			successHandler := createFlowConsumerSuccessHandler(flowName, broker, ID, rawInputEvent, flows, outputVarNames, vars, &lock, &allConsumerDeclared, outputCompleted)
			successHandlers[rawInputEvent] = successHandler
			errorHandler := createFlowConsumerErrorHandler(flowName, outputCompleted)
			log.Printf("[INFO: flow.%s] Consuming from %s", flowName, inputEvent)
			broker.Subscribe(inputEvent, successHandler, errorHandler)
			allConsumerDeclared.Done()
		}
		// set default vars (`var -> output` scenario)
		setFlowDefaultVars(flowName, flows, outputVarNames, vars, &lock, outputCompleted)
		allConsumerDeclared.Wait()
		// publishServiceOutput default vars
		lock.RLock()
		publishFlowVar(flowName, broker, ID, flows, outputVarNames, getMapKeys(vars), vars)
		lock.RUnlock()
		// publishServiceOutput preset vars
		executed := map[string]bool{}
		for presetVarName, presetValue := range presetVars {
			for _, rawInputEvent := range flows.GetInputEventByVarName(presetVarName) {
				if !executed[rawInputEvent] {
					handler := successHandlers[rawInputEvent]
					handler(servicedata.Package{ID: ID, Data: presetValue})
					executed[rawInputEvent] = true
				}
			}
		}
		<-outputCompleted
		lock.RLock()
		outputs := createFlowOutputs(flowName, outputVarNames, vars)
		lock.RUnlock()
		log.Printf("[INFO: flow.%s] Internal flow `%s` ended. Outputs are: `%#v`", flowName, ID, outputs)
		// unsubscribe
		for _, rawInputEvent := range rawInputEvents {
			inputEvent := fmt.Sprintf("%s.%s", ID, rawInputEvent)
			err := broker.Unsubscribe(inputEvent)
			if err != nil {
				return outputs, err
			}
		}
		return outputs, err
	}
}

func setFlowDefaultVars(flowName string, flows FlowEvents, outputVarNames []string, vars Dictionary, lock *sync.RWMutex, outputCompleted chan bool) {
	defaultVars := getFlowDefaultVars(flows)
	for varName, value := range defaultVars {
		log.Printf("[INFO: flow.%s] Internally set `%s` into: `%#v`", flowName, varName, value)
		err := vars.Set(varName, value)
		if err != nil {
			log.Printf("[ERROR: flow.%s] Error setting `%s`: %s", flowName, varName, err)
		}
	}
	lock.RLock()
	completed := vars.HasAll(outputVarNames)
	lock.RUnlock()
	if completed {
		outputCompleted <- true
	}
}

func createFlowConsumerSuccessHandler(flowName string, broker msgbroker.CommonBroker, ID, rawInputEvent string, flows FlowEvents, outputVarNames []string, vars Dictionary, lock *sync.RWMutex, allConsumerDeclared *sync.WaitGroup, outputCompleted chan bool) func(servicedata.Package) {
	return func(pkg servicedata.Package) { // consume success
		inputEvent := fmt.Sprintf("%s.%s", ID, rawInputEvent)
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
			// this will be executed on: input -> var -> output scenario
			publishedData := pkg.Data
			if varFlow.UseValue {
				publishedData = varFlow.Value
			}
			if varFlow.UseFunction && varFlow.Function != nil {
				publishedData = varFlow.Function(publishedData)
			}
			log.Printf("[INFO: flow.%s] Set `%s` into: `%#v`", flowName, varName, publishedData)
			lock.Lock()
			err := vars.Set(varName, publishedData)
			if err != nil {
				log.Printf("[ERROR: flow.%s] Error setting `%s`: %s", flowName, varName, err)
			}
			lock.Unlock()
			lock.RLock()
			publishFlowVar(flowName, broker, ID, flows, outputVarNames, []string{varName}, vars)
			lock.RUnlock()
		}
		lock.RLock()
		completed := vars.HasAll(outputVarNames)
		lock.RUnlock()
		if completed {
			outputCompleted <- true
		}
	}
}

func createFlowConsumerErrorHandler(flowName string, outputCompleted chan bool) func(error) {
	return func(err error) { // consume error
		log.Printf("[ERROR: flow.%s] Error: %s", flowName, err)
		outputCompleted <- true
	}
}

func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}
	return keys
}

func createFlowOutputs(flowName string, outputVarNames []string, vars Dictionary) Dictionary {
	outputs := make(Dictionary)
	for _, outputName := range outputVarNames {
		err := outputs.Set(outputName, vars.Get(outputName))
		if err != nil {
			log.Printf("[ERROR: flow.%s] Error creating output: %s", flowName, err)
		}
	}
	return outputs
}

func getFlowRawInputEvents(flows FlowEvents, inputVarNames []string) []string {
	candidates := flows.GetInputEvents()
	var rawInputEvents []string
	for _, candidate := range candidates {
		rawInputEvents = append(rawInputEvents, candidate)
	}
	return rawInputEvents
}

func getFlowDefaultVars(flows FlowEvents) map[string]interface{} {
	defaultVars := make(map[string]interface{})
	// determine candidates from flows
	candidates := flows.GetVarNamesByInputEvent("")
	for _, candidate := range candidates {
		candidatePass := false
		var value interface{}
		var function func(interface{}) interface{}
		useFunction := false
		for _, flow := range flows {
			if flow.VarName == candidate && flow.InputEvent == "" && flow.UseValue {
				value = flow.Value
				function = flow.Function
				useFunction = flow.UseFunction
				candidatePass = true
			} else if isSubVarOf(flow.VarName, candidate) {
				candidatePass = false
				break
			}
		}
		if candidatePass {
			if useFunction && function != nil {
				value = function(value)
			}
			defaultVars[candidate] = value
		}
	}
	return defaultVars
}

func publishFlowVar(flowName string, broker msgbroker.CommonBroker, ID string, flows FlowEvents, exceptions, publishedVarNames []string, vars Dictionary) {
	for _, publishedVarName := range publishedVarNames {
		// if var is part of exceptions, ignore it. CommonBroker will do the job
		if IsStringInArray(publishedVarName, exceptions) {
			continue
		}
		// varNames contains publishedVarName and all it's sub variable's names
		varNames := []string{publishedVarName}
		for _, flow := range flows {
			varName := flow.VarName
			if varName != publishedVarName && isSubVarOf(publishedVarName, varName) {
				varNames = append(varNames, varName)
			}
		}
		// for every varNames, get it's outputEvent and publishServiceOutput
		for _, varName := range varNames {
			for _, rawOutputEvent := range flows.GetOutputEventByVarNames(varName) {
				varValue := vars.Get(varName)
				outputEvent := fmt.Sprintf("%s.%s", ID, rawOutputEvent)
				err := Publish("flow", flowName, broker, ID, outputEvent, varValue)
				if err != nil {
					log.Printf("[INFO: flow.%s] Error while publishing into `%s`: %s", flowName, outputEvent, err)
				}
			}
		}
	}
}

// isSubVarOf is current sub of candidate
func isSubVarOf(candidate, current string) bool {
	return strings.Index(current, fmt.Sprintf("%s.", candidate)) == 0
}
