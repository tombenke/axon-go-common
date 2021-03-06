package nats

import (
	nats "github.com/nats-io/nats.go"
	messenger "github.com/tombenke/axon-go-common/messenger"
	"time"
)

// Publish `msg` message to the `subject` topic
func (m connections) Publish(subject string, msg []byte) error {
	subj := subject

	if err := m.nc.Publish(subj, msg); err != nil {
		m.logger.Fatalf("Messenger error: '%s'", err.Error())
		return err
	}
	m.nc.Flush()

	m.logger.Debugf("Messenger published message to '%s'", subj)
	return nil
}

// Subscriber structure holds the subscriptions
// and provides methods that implements the generic Subscriber
type Subscriber struct {
	Subscription *nats.Subscription
}

// Unsubscribe unsubscribes the `Subscriber` from the topic
func (subs Subscriber) Unsubscribe() error {
	return subs.Subscription.Unsubscribe()
}

// newSubscriber creates a new NATS Subscriber
func newSubscriber(subscription *nats.Subscription) messenger.Subscriber {
	var subscriber = Subscriber{Subscription: subscription}
	return subscriber
}

// Subscribe subscribes to the `subject` topic, and calls the `cb` call-back function with the inbound messages
func (m connections) Subscribe(subject string, cb func([]byte)) messenger.Subscriber {
	subscription, err := m.nc.Subscribe(subject, func(msg *nats.Msg) {
		cb(msg.Data)
	})
	if err != nil {
		panic(err)
	}
	m.nc.Flush()
	return newSubscriber(subscription)
}

// ChanSubscribe subscribes to the `subject` topic, and sends the inbound messages into the `ch` channel
// You should not close the channel until sub.Unsubscribe() has been called.
func (m connections) ChanSubscribe(subject string, ch chan []byte) messenger.Subscriber {
	subscription, err := m.nc.Subscribe(subject, func(msg *nats.Msg) {
		m.logger.Debugf("Messenger received message from '%s'", subject)
		ch <- msg.Data
	})
	if err != nil {
		panic(err)
	}
	m.nc.Flush()
	return newSubscriber(subscription)
}

// Request `msg` message through the `subject` topic and expects a response until `timeout`.
func (m connections) Request(subject string, msg []byte, timeout time.Duration) ([]byte, error) {
	subj := subject

	m.logger.Debugf("Messenger sends request through '%s'", subj)
	resp, err := m.nc.Request(subj, msg, timeout)
	if err != nil {
		m.logger.Error(err)
		return nil, err
	}

	m.logger.Debugf("Messenger got response '%s'", resp.Data)
	return resp.Data, err
}

// Subscribe to the `subject` topic, and calls the `service` call-back function with the inbound messages,
// then respond with the return value of the `service` function through the `Reply` subject.
func (m connections) Response(subject string, service func([]byte) ([]byte, error)) {
	_, err := m.nc.Subscribe(subject, func(msg *nats.Msg) {
		resp, err := service(msg.Data)
		if err != nil {
			if err := m.nc.Publish(msg.Reply, []byte(err.Error())); err != nil {
				panic(err)
			}
		} else {
			if err := m.nc.Publish(msg.Reply, resp); err != nil {
				panic(err)
			}
		}
	})
	if err != nil {
		panic(err)
	}
	m.nc.Flush()
}
