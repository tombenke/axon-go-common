package orchestra

import (
	"encoding/json"
	"fmt"
	"github.com/tombenke/axon-go-common/msgs"
	"github.com/tombenke/axon-go-common/msgs/common"
	"time"
)

const (
	// EPNStatusTypeName is the printable name of the `EPNStatus` message-type
	EPNStatusTypeName = "orchestra/EPNStatus"
)

func init() {
	msgs.RegisterMessageType(EPNStatusTypeName, []msgs.Representation{msgs.JSONRepresentation}, func() msgs.Message {
		return NewEPNStatusMessage(EPNStatusBody{})
	})
}

// EPNStatus represents the structure of the `status-report` message that the actor is usually sent
// to the orchestrator as a response to the `status-request` message.
// The `Body.Data` property holds the name of the actor that sends the response.
type EPNStatus struct {
	Header common.Header
	Body   EPNStatusBody
}

// EPNStatusBody represents the entire Event Processing Network that is made of actor nodes and channels.
type EPNStatusBody struct {
	// Actors holds the list of actors
	Actors []Actor
}

// Actor represents the status of one actor node
type Actor struct {
	Name         string
	ResponseTime time.Duration
}

// GetType returns with the printable name of the `EPNStatus` message-type
func (msg *EPNStatus) GetType() string {
	return EPNStatusTypeName
}

// Encode returns with the `EPNStatus` message content in a representation format selected by `representation`
func (msg *EPNStatus) Encode(representation msgs.Representation) (results []byte) {
	switch representation {
	case msgs.JSONRepresentation:
		var err error
		results, err = json.Marshal(*msg)
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Errorf("Encode error: unknown representational format '%s'", representation))
	}
	return results
}

// Decode parses the `content` using the selected `representation` format
func (msg *EPNStatus) Decode(representation msgs.Representation, content []byte) error {
	switch representation {
	case msgs.JSONRepresentation:
		return json.Unmarshal(content, msg)
	default:
		panic(fmt.Errorf("Decode error: unknown representational format '%s'", representation))
	}
}

// JSON returns with the `EPNStatus` message content in JSON representation format
func (msg *EPNStatus) JSON() []byte {
	jsonBytes, err := json.Marshal(*msg)
	if err != nil {
		panic(err)
	}
	return jsonBytes
}

// EPNStatus returns with the `EPNStatus` message content in JSON format string
func (msg *EPNStatus) String() string {
	jsonBytes, err := json.Marshal(*msg)
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

// ParseJSON parses the JSON representation of a `EPNStatus` messages from the `jsonBytes` argument.
func (msg *EPNStatus) ParseJSON(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, msg)
}

// NewEPNStatusMessage returns with a new `EPNStatus` message. The header will contain the current time in `Nanoseconds` precision.
func NewEPNStatusMessage(body EPNStatusBody) msgs.Message {
	return NewEPNStatusMessageAt(body, time.Now().UnixNano(), "ns")
}

// NewEPNStatusMessageAt returns with a new `EPNStatus` message. The header will contain the `at` time in `withPrecision` precision.
func NewEPNStatusMessageAt(body EPNStatusBody, at int64, withPrecision common.TimePrecision) msgs.Message {
	var msg EPNStatus
	msg.Header = common.NewHeaderAt(at, withPrecision)
	msg.Body = body
	return &msg
}
