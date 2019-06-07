package gen

// InputEvent definition
type InputEvent struct {
	EventName string
	VarName   string
}

// NewInputEvent create new InputEvent
func NewInputEvent(eventName, varName string) InputEvent {
	return InputEvent{EventName: eventName, VarName: varName}
}
