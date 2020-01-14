package axon

import (
    "testing"
    "fmt"
	"github.com/nats-io/nats.go"
)

func TestSetupConnoptions(t *testing.T) {
	opts := []nats.Option{nats.Name("natsTest")}
	opts = setupConnOptions(opts)

    for i, v := range(opts) {
        fmt.Printf("opts: %v, %v", i, v)
    }

    if l := len(opts); l != 6 {
        t.Error("setupConnoptions should return with 6 options")
    }
}
