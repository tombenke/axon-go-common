// The axon package implements agents for IoT applications.

/*
The axon package implements agents for IoT applications.

About

Axon is a set of independent components, that can be written in any programming languages,
which has a library to access https://nats.io/.

The components are event driven agents that either consume and/or produce messages through NATS.
These agents use NATS subjects for communicating with each others.

The structure of the messages a given kind of agent is able to consume,
or produces depents on the given agent, as well as its behavior.

The axon package contains a `common` module, that provides generic functions for the agents,
e.g. connecting to the NATS server, etc.

The package also contain a set of predefined agents, such as `axon-cron`, `axon-debug`,
that are compiled and can be executed as standalone applications.

From a given perspective, axon is similar to the Node-RED (https://nodered.org/)
in the meaning that its agents work similarly like the Node-RED components.
There are three fundamental differences:

1. the axon agents' inputs and outputs are NATS subjects, or channels,

2. the agents can be written in any language,

3. the agents can run on different machines and in any number of instances.

For more detailed description, read the README file of the project (https://github.com/tombenke/axon-go).
*/
package axon

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// ConnectToNats connects to the NATS server and returns with a `*nats.Conn` that can be used for further operations.
func ConnectToNats(urls string, userCreds string, natsName string) (*nats.Conn, error) {
	// Connect Options.
	opts := []nats.Option{nats.Name(natsName)}
	opts = setupConnOptions(opts)

	// Use UserCredentials
	if userCreds != "" {
		opts = append(opts, nats.UserCredentials(userCreds))
	}

    return nats.Connect(urls, opts...)
}

// setupOptions extends the options array with default configuration parameters
// that are useful to connect to the NATS server.
func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectHandler(func(nc *nats.Conn) {
		log.Printf("Disconnected: will attempt reconnects for %.0fm", totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}

// Check if there is an error, and print it with the `prefix`.
func Check(prefix string, err error) {
    if err != nil {
        log.Printf("%s%s", prefix, err)
    }
}

func CheckFatal(err error) {
    if err != nil {
        log.Fatal(err)
        panic(err)
    }
}

func SendMessage(nc *nats.Conn, subject string, msg []byte) error {
    subj := subject

    nc.Publish(subj, msg)
    nc.Flush()

    if err := nc.LastError(); err != nil {
        log.Fatal(err)
    } else {
        log.Printf("Published [%s] : '%s'\n", subj, msg)
    }
    return nil
}

// Move to axon-go-common
func NowAsUnixWithPrecision(precision string) int64 {
    nowNs := time.Now().UnixNano()
    switch precision {
    case "ns":
        return nowNs
    case "u", "us":
        return nowNs / 1e3
    case "ms":
        return nowNs / 1e6
    case "s":
        return nowNs / 1e9
    }
    return nowNs
}

