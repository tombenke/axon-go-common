package orchestra

import (
	"encoding/json"
	"fmt"
	"github.com/tombenke/axon-go-common/msgs"
	"github.com/tombenke/axon-go-common/msgs/common"
	"time"
)

const (
	// StatusReportTypeName is the printable name of the `StatusReport` message-type
	StatusReportTypeName = "orchestra/StatusReport"
)

func init() {
	msgs.RegisterMessageType(StatusReportTypeName, []msgs.Representation{msgs.JSONRepresentation}, func() msgs.Message {
		return NewStatusReportMessage(StatusReportBody{})
	})
}

// StatusReport represents the structure of the `status-report` message that the actor is usually sent
// to the orchestrator as a response to the `status-request` message.
// The `Body.Data` property holds the name of the actor that sends the response.
type StatusReport struct {
	Header common.Header
	Body   StatusReportBody
}

// StatusReportBody represents the internal structure of the message the node sends as a response
// through the 'status-report' channel to the orchestrator. It holds the detailed description of the
// node, incl. the main characteristics of its ports.
type StatusReportBody struct {

	// Name is the name of the node. It should be unique in a specific network
	Name string

	// Type is the symbolic name of the node type, that refers to how the node is working.
	Type string

	// Ports holds the I/O port definitions
	Ports Ports

	// Synchronization is a flag. If it is `true` the Node is working in syncronized mode,
	// otherwise it uses no synchronization protocol.
	Synchronization bool

	// SpecsURL holds an URL to the base-path of the detailed specification of the Node.
	SpecsURL string
}

type Ports struct {

	// Inputs is a list of input-type port descriptors
	Inputs []Port

	// Outputs is a list of output-type port descriptors
	Outputs []Port
}

// IO defines the properties of a generic I/O port
type Port struct {
	// Name is the name of the port
	Name string

	// Type is the message-type the port uses for transfer
	Type string

	// Representation is the message representation format used for transfer
	Representation string

	// Channel is the representation of a communication channel of a port of a node
	Channel Channel
}

// Channel represents a messaging subject, that the ports use for communication
type Channel struct {
	// Name is the name of the messaging subject
	Name string

	// Type is the type of the messaging subject. Valid values: TOPIC, WORKER, RPC, CHANNEL.
	Type string
}

// GetType returns with the printable name of the `StatusReport` message-type
func (msg *StatusReport) GetType() string {
	return StatusReportTypeName
}

// Encode returns with the `StatusReport` message content in a representation format selected by `representation`
func (msg *StatusReport) Encode(representation msgs.Representation) (results []byte) {
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
func (msg *StatusReport) Decode(representation msgs.Representation, content []byte) error {
	switch representation {
	case msgs.JSONRepresentation:
		return json.Unmarshal(content, msg)
	default:
		panic(fmt.Errorf("Decode error: unknown representational format '%s'", representation))
	}
}

// JSON returns with the `StatusReport` message content in JSON representation format
func (msg *StatusReport) JSON() []byte {
	jsonBytes, err := json.Marshal(*msg)
	if err != nil {
		panic(err)
	}
	return jsonBytes
}

// StatusReport returns with the `StatusReport` message content in JSON format string
func (msg *StatusReport) String() string {
	jsonBytes, err := json.Marshal(*msg)
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

// ParseJSON parses the JSON representation of a `StatusReport` messages from the `jsonBytes` argument.
func (msg *StatusReport) ParseJSON(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, msg)
}

// NewStatusReportMessage returns with a new `StatusReport` message. The header will contain the current time in `Nanoseconds` precision.
func NewStatusReportMessage(body StatusReportBody) msgs.Message {
	return NewStatusReportMessageAt(body, time.Now().UnixNano(), "ns")
}

// NewStatusReportMessageAt returns with a new `StatusReport` message. The header will contain the `at` time in `withPrecision` precision.
func NewStatusReportMessageAt(body StatusReportBody, at int64, withPrecision common.TimePrecision) msgs.Message {
	var msg StatusReport
	msg.Header = common.NewHeaderAt(at, withPrecision)
	msg.Body = body
	return &msg
}
