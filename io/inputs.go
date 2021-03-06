package io

import (
	"fmt"
	"github.com/tombenke/axon-go-common/config"
	"github.com/tombenke/axon-go-common/msgs"
	"sync"
)

// Input holds the data of an input port of the actor
type Input struct {
	IO
	DefaultMessage msgs.Message
}

// Inputs holds a map of the the input ports of the actor. The key is the name of the port.
type Inputs struct {
	RW  sync.RWMutex
	Map map[string]Input
}

////type Inputs map[string]Input

// GetMessage returns the last message received via the input port selected by the `name` parameter
func (inputs *Inputs) GetMessage(name string) msgs.Message {

	(*inputs).RW.RLock()
	defer (*inputs).RW.RUnlock()

	if input, ok := inputs.Map[name]; ok {
		return input.Message
	}
	errorMessage := fmt.Sprintf("There is no input port named to '%s'", name)
	panic(errorMessage)
}

// SetMessage sets the message that received via the input channel to the port selected by the `name` parameter
func (inputs *Inputs) SetMessage(name string, inMsg msgs.Message) {
	(*inputs).RW.Lock()
	defer (*inputs).RW.Unlock()

	if _, ok := (*inputs).Map[name]; !ok {
		errorMessage := fmt.Sprintf("'%s' port does not exist, so can not set message to it.", name)
		panic(errorMessage)
	}

	inMsgType := inMsg.GetType()
	portMsgType := string((*inputs).Map[name].Type)
	if inMsgType != portMsgType {
		errorMessage := fmt.Sprintf("'%s' message-type mismatch to port's '%s' message-type.", inMsgType, portMsgType)
		panic(errorMessage)
	}

	(*inputs).Map[name] = Input{
		IO: IO{
			Name:           name,
			Type:           inMsgType,
			Representation: (*inputs).Map[name].Representation,
			Channel:        (*inputs).Map[name].Channel,
			Message:        inMsg,
		},
		DefaultMessage: (*inputs).Map[name].DefaultMessage,
	}
}

// NewInputs creates a new Inputs map based on the config parameters
func NewInputs(inputsCfg config.Inputs) *Inputs {
	inputs := Inputs{
		RW:  *new(sync.RWMutex),
		Map: make(map[string]Input),
	}

	inputs.RW.Lock()
	defer inputs.RW.Unlock()

	for _, in := range inputsCfg {
		inputs.Map[in.Name] = NewInput(in.IO.Name, in.IO.Type, msgs.Representation(in.IO.Representation), in.IO.Channel, NewDefaultMessage(in.Type, in.Default))
	}
	return &inputs
}

// NewInput create a new Input message instance according to the arguments
func NewInput(Name string, Type string, Repr msgs.Representation, Chan string, Default msgs.Message) Input {

	// Validates if the message-type is registered
	if !msgs.IsMessageTypeRegistered(Type) {
		errorString := fmt.Sprintf("The '%s' message type has not been registered!", Type)
		panic(errorString)
	}

	// Validates if the representation format is supported
	if !msgs.DoesMessageTypeImplementsRepresentation(Type, Repr) {
		errorString := fmt.Sprintf("'%s' message-type does not implement codec for '%s' representation format", Type, Repr)
		panic(errorString)
	}

	return Input{IO: IO{Name: Name, Type: Type, Representation: Repr, Channel: Chan, Message: Default}, DefaultMessage: Default}
}

// NewDefaultMessage create a new default message of `Type` from the `newDefault` string value.
// If `newDefault` is empty, then creates the built-in default message defined to the specific Type.
func NewDefaultMessage(Type string, newDefault string) msgs.Message {
	// Determines the default value
	// In case the default value is empty, then use the original one defined by the message-type itself
	Default := msgs.GetDefaultMessageByType(Type)
	if newDefault != "" {
		// The default config value is not empty, so it should be a valid message in JSON format
		err := Default.Decode(msgs.JSONRepresentation, []byte(newDefault))
		if err != nil {
			panic(err)
		}
	}

	if Default == nil {
		panic("Wrong NewDefaultMessage")
	}

	return Default
}
