package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// invalid input strings
var invalidIns []string = []string{
	"",              // no name string
	"|",             // empty name string
	"||",            // empty name string
	"|channel||",    // empty name string
	"||0.1",         // empty name string
	"|channel||0.1", // empty name string
	"name||||||",    // Wrong number of arguments
	"name|||",       // Wrong number of arguments
}

type validIn struct {
	Arg      string
	Expected In
}

var validIns []validIn = []validIn{
	validIn{"name", In{IO{"name", DefaultType, DefaultRepresentation, ""}, ""}},                                                 // name only
	validIn{"name||||0.1", In{IO{"name", DefaultType, DefaultRepresentation, ""}, "0.1"}},                                       // name and default value
	validIn{"name||||0.1", In{IO{"name", DefaultType, DefaultRepresentation, ""}, "0.1"}},                                       // name and default value
	validIn{"name|channel|||", In{IO{"name", DefaultType, DefaultRepresentation, "channel"}, ""}},                               // channel and name
	validIn{"name|channel|||false", In{IO{"name", DefaultType, DefaultRepresentation, "channel"}, "false"}},                     // channel and name
	validIn{"name|channel|base/Bool|application/json|true", In{IO{"name", "base/Bool", "application/json", "channel"}, "true"}}, // full
}

// Test input args
func TestParseInArgs(t *testing.T) {
	assert := assert.New(t)

	// Test valid cases
	for _, i := range validIns {
		assert.Equal(i.Expected, parseIn(i.Arg))
	}

	// Test invalid cases
	for _, i := range invalidIns {
		assert.Panics(
			func() {
				parseIn(i)
			},
			"It should panic!",
		)
	}
}

// invalid output strings
var invalidOuts []string = []string{
	"",          // no name string
	"|",         // empty name string
	"|channel",  // empty name string
	"name||",    // Wrong number of arguments
	"name|||||", // Wrong number of arguments
}

type validOut struct {
	Arg      string
	Expected Out
}

var validOuts []validOut = []validOut{
	validOut{"name", Out{IO{"name", DefaultType, DefaultRepresentation, ""}}},
	validOut{"name|", Out{IO{"name", DefaultType, DefaultRepresentation, ""}}},
	validOut{"name|channel|base/Bool|application/json", Out{IO{"name", "base/Bool", "application/json", "channel"}}},
}

// Test output args
func TestParseOutArgs(t *testing.T) {
	assert := assert.New(t)

	// Test valid cases
	for _, i := range validOuts {
		assert.Equal(parseOut(i.Arg), i.Expected)
	}

	// Test invalid cases
	for _, i := range invalidOuts {
		assert.Panics(
			func() {
				parseOut(i)
			},
			"It should panic!",
		)
	}
}

func TestSetIn(t *testing.T) {
	inputs := &Inputs{}
	assert.Nil(t, inputs.Set("name|channel|base/Bool|application/json|true"))
	assert.Nil(t, inputs.Set("name2|channel2|base/Any|application/json|{}"))
	assert.Nil(t, inputs.Set("name|channelx|base/Bytes|text/plain|")) // Overwrites 'name' !!!!
	assert.Nil(t, inputs.Set(`name3|channel3|base/Float|application/json|{"Body":{"Data":42.}}`))

	expected := Inputs{
		In{IO{"name", "base/Bytes", "text/plain", "channelx"}, ""},
		In{IO{"name2", "base/Any", "application/json", "channel2"}, "{}"},
		In{IO{"name3", "base/Float", "application/json", "channel3"}, `{"Body":{"Data":42.}}`},
	}
	assert.Equal(t, expected, *inputs)
}

func TestSetOut(t *testing.T) {
	outputs := &Outputs{}
	assert.Nil(t, outputs.Set("name|channel|base/Bool|application/json"))
	assert.Nil(t, outputs.Set("name2|channel2|base/Any|application/json"))
	assert.Nil(t, outputs.Set("name|channelx|base/Bytes|text/plain")) // Overwrites 'name' !!!!
	assert.Nil(t, outputs.Set("name3|channel3|base/Float|application/json"))

	expected := Outputs{
		Out{IO{"name", "base/Bytes", "text/plain", "channelx"}},
		Out{IO{"name2", "base/Any", "application/json", "channel2"}},
		Out{IO{"name3", "base/Float", "application/json", "channel3"}},
	}
	assert.Equal(t, expected, *outputs)
}
