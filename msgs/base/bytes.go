package base

import (
	"encoding/json"
	"fmt"
	"github.com/tombenke/axon-go-common/msgs"
)

const (
	// BytesTypeName is the printable name of the `Bytes` message-type
	BytesTypeName = "base/Bytes"
)

func init() {
	msgs.RegisterMessageType(BytesTypeName, []msgs.Representation{msgs.TextRepresentation}, func() msgs.Message {
		return NewBytesMessage([]byte{})
	})
}

// Bytes represents the structure of a generic message that is actually a plain byte array
type Bytes []byte

// GetType returns with the printable name of the `Bytes` message-type
func (msg *Bytes) GetType() string {
	return BytesTypeName
}

// Encode returns with the `Bytes` message content in a representation format selected by `representation`
func (msg *Bytes) Encode(representation msgs.Representation) (results []byte) {
	switch representation {
	case msgs.TextRepresentation:
		return results
	default:
		panic(fmt.Errorf("Encode error: unknown representational format '%s'", representation))
	}
	return results
}

// Decode parses the `content` using the selected `representation` format
func (msg *Bytes) Decode(representation msgs.Representation, content []byte) error {
	switch representation {
	case msgs.TextRepresentation:
		return nil
	default:
		panic(fmt.Errorf("Decode error: unknown representational format '%s'", representation))
	}
}

// String returns with the `Bytes` message content in JSON format string
func (msg *Bytes) String() string {
	return string(*msg)
}

// JSON returns with the `Any` message content in JSON representation format
func (msg *Bytes) JSON() []byte {
	jsonBytes, err := json.Marshal(string(*msg))
	if err != nil {
		panic(err)
	}
	return jsonBytes
}

// ParseJSON parses the JSON representation of a `Any` messages from the `jsonBytes` argument.
func (msg *Bytes) ParseJSON(jsonBytes []byte) error {
	if len([]byte{}) == 0 {
		(*msg) = Bytes{}
		return nil
	}
	return json.Unmarshal(jsonBytes, msg)
	//return json.Unmarshal([]byte("\""+string(jsonBytes)+"\""), msg)
}

// NewBytesMessage returns with a new `Bytes` message. The header will contain the current time in `Nanoseconds` precision.
func NewBytesMessage(data []byte) msgs.Message {
	var msg Bytes = data
	return &msg
}
