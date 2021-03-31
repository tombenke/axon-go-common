package base

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/tombenke/axon-go-common/msgs"
	"testing"
)

func TestBytesGetType(t *testing.T) {
	data := []byte{}
	assert.Equal(t, NewBytesMessage(data).GetType(), BytesTypeName)
}

func TestBytesMessage(t *testing.T) {
	data := []byte{}
	m := NewBytesMessage(data)
	var n Bytes
	err := n.ParseJSON(m.JSON())
	assert.Nil(t, err)
	err = n.ParseJSON([]byte(m.String()))
	assert.Nil(t, err)
	assert.Equal(t, m, &n)
}

func TestBytesMessageCodec(t *testing.T) {
	data := []byte{}
	m := NewBytesMessage(data)
	n := Bytes{}
	err := n.Decode(msgs.TextRepresentation, m.Encode(msgs.TextRepresentation))
	assert.Nil(t, err)
	assert.Equal(t, m, &n)
	err = n.Decode(msgs.OctetstreamRepresentation, m.Encode(msgs.OctetstreamRepresentation))
	assert.Nil(t, err)
	assert.Equal(t, m, &n)
}

func TestBytesMessageCodecPanic(t *testing.T) {
	data := []byte{}
	m := NewBytesMessage(data)
	n := Bytes{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, r, errors.New("Decode error: unknown representational format 'wrong-representation'"))
			}
		}()
		err := n.Decode(msgs.Representation("wrong-representation"), m.Encode(msgs.TextRepresentation))
		assert.Nil(t, err)
	}()
	func() {
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, r, errors.New("Encode error: unknown representational format 'wrong-representation'"))
			}
		}()
		err := n.Decode(msgs.TextRepresentation, m.Encode(msgs.Representation("wrong-representation")))
		assert.Nil(t, err)
	}()
}
