package service

import (
	"encoding/json"
	"fmt"
	"log"
	nats "github.com/nats-io/nats.go"
)

// WrappedFunction any function wrapped for ConsumeAndPublish
type WrappedFunction = func(inputs Dictionary) Dictionary

// CommonService single configuration
type CommonService struct {
	Input    []IO
	Output   []IO
	Function WrappedFunction
}

// NewCommonService create new singleConfig for service
func NewCommonService(serviceName string, methodName string, inputs []string, outputs []string, wrappedFunction WrappedFunction) CommonService {
	// get inputConfig
	var inputConfig []IO
	for _, inputName := range inputs {
		eventName := fmt.Sprintf("srvc.%s.%s.in.%s", serviceName, methodName, inputName)
		inputConfig = append(inputConfig, IO{VarName: inputName, EventName: eventName})
	}
	// get outputConfig
	var outputConfig []IO
	for _, outputName := range outputs {
		eventName := fmt.Sprintf("srvc.%s.%s.in.%s", serviceName, methodName, outputName)
		outputConfig = append(outputConfig, IO{VarName: outputName, EventName: eventName})
	}
	// return config
	return CommonService{
		Input:    inputConfig,
		Output:   outputConfig,
		Function: wrappedFunction,
	}
}

// Services configuration
type Services = map[string]CommonService

// ConsumeAndPublish consume from queue and Publish
func (services Services)ConsumeAndPublish() {
	natsURL := GetNatsURL()
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Print(err)
		return
	}
	for _, service := range services {
		inputConfig := service.Input
		outputConfig := service.Output
		wrappedFunction := service.Function
		consumeAndPublishSingle(nc, inputConfig, outputConfig, wrappedFunction)
	}
}

func consumeAndPublishSingle(nc *nats.Conn, inputConfig []IO, outputConfig []IO, wrappedFunction WrappedFunction) {
	// allInputs
	allInputs := make(map[string]Dictionary)
	rawInputEventNames := GetUniqueEventNames(inputConfig)
	inputVarNames := GetUniqueVarNames(inputConfig)
	for _, rawInputEventName := range rawInputEventNames {
		inputEventName := fmt.Sprintf("*.%s", rawInputEventName)
		log.Printf("[INFO] Consume from `%s`", inputEventName)
		nc.Subscribe(inputEventName, func(m *nats.Msg) {
			var pkg Package
			JSONByte := m.Data
			log.Printf("[INFO] Get message from `%s`: %s", inputEventName, string(JSONByte))
			err := json.Unmarshal(JSONByte, &pkg)
			if err != nil {
				log.Printf("[ERROR] %s: %s", inputEventName, err)
				return
			}
			// prepare allInputs
			ID := pkg.ID
			data := pkg.Data
			if _, exists := allInputs[ID]; !exists {
				allInputs[ID] = Dictionary{}
			}
			// populate allInputs[ID] with eventInputNames and data
			eventInputNames := GetEventVarNames(inputConfig, rawInputEventName)
			for _, inputVarName := range eventInputNames {
				allInputs[ID][inputVarName] = data
			}
			inputs := allInputs[ID]
			log.Printf("[INFO] Inputs for %s: %#v", ID, inputs)
			// execute wrapper
			if isInputComplete(inputVarNames, inputs) {
				log.Printf("[INFO] Inputs for %s completed", ID)
				outputs := wrappedFunction(inputs)
				publish(nc, ID, outputConfig, outputs)
			}
		})
	}
}

func isInputComplete(inputVarNames []string, inputs Dictionary) bool {
	for _, inputVarName := range inputVarNames {
		if _, exists := inputs[inputVarName]; !exists {
			return false
		}
	}
	return true
}

func publish(nc *nats.Conn, ID string, outputConfig []IO, outputs Dictionary) {
	outputVarNames := GetUniqueVarNames(outputConfig)
	for _, outputVarName := range outputVarNames {
		data := outputs[outputVarName]
		pkg := Package{ID: ID, Data: data}
		rawOutputEventNames := GetVarEventNames(outputConfig, outputVarName)
		for _, rawOutputEventName := range rawOutputEventNames {
			publishPkg(nc, ID, rawOutputEventName, pkg)
		}
	}
}

func publishPkg(nc *nats.Conn, ID string, rawOutputEventName string, pkg Package) {
	JSONByte, err := json.Marshal(&pkg)
	if err != nil {
		log.Printf("[ERROR] %s: %s", ID, err)
		return
	}
	eventName := fmt.Sprintf("%s.%s", ID, rawOutputEventName)
	nc.Publish(eventName, JSONByte)
	log.Printf("[INFO] Publish into `%s`: `%#v`", eventName, pkg)
}
