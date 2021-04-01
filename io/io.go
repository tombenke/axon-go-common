package io

import (
	"github.com/tombenke/axon-go-common/msgs"
)

// IO is a generic port descriptor structure,
// that contains the common properties of both the input and output ports
type IO struct {
	// Name The name of the port that uniquely identifies it
	Name string
	// Type The `Type` of the messages it receives or emits,
	Type string
	// Representation The `Representation` format the input ports decodes from
	// and the output ports encodes to the internal representation of the messages
	Representation msgs.Representation
	// Channel The name of the `Channel` the input port receives, and the output ports sends the messages
	Channel string
	// Message The actual message the io port holds
	Message msgs.Message
}

// Handler is an interface for both the input and output type ports
type Handler interface {
	GetInputMessage(string) (msgs.Message, error)
	SetOutputMessage(string, msgs.Message) error
}
