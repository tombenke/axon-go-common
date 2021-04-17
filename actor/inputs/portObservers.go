package inputs

import (
	"github.com/sirupsen/logrus"
	"github.com/tombenke/axon-go-common/io"
	"github.com/tombenke/axon-go-common/messenger"
	"github.com/tombenke/axon-go-common/msgs"
	"sync"
)

// startInPortsObservers starts one message observer for every port,
// and returns with the number of observers started.
func startInPortsObservers(inputs *io.Inputs, inputsMuxCh chan io.Input, doneCh chan bool, wg *sync.WaitGroup, m messenger.Messenger, logger *logrus.Logger) {
	for p := range (*inputs).Map {
		if (*inputs).Map[p].Channel != "" {
			startedCh := newPortObserver((*inputs).Map[p], inputsMuxCh, doneCh, wg, m, logger)
			<-startedCh
		}
	}
}

// newPortObserver subscribes to an input channel with a go routine that observes the incoming messages.
// When a message arrives through the channel, the go routine forwards that through the `inCh` towards the aggregator.
// The newPortObserver creates and returns with the `inCh` channel that the aggregator can consume.
func newPortObserver(input io.Input, inputsMuxCh chan io.Input, doneCh chan bool, wg *sync.WaitGroup, m messenger.Messenger, logger *logrus.Logger) chan interface{} {
	inMsgCh := make(chan []byte)
	logger.Debugf("Receiver's '%s' port observer subscribe to '%s' channel", input.Name, input.Channel)
	inMsgSubs := m.ChanSubscribe(input.Channel, inMsgCh)
	startedCh := make(chan interface{})

	wg.Add(1)
	go func() {
		logger.Debugf("Receiver's '%s' port observer started", input.Name)
		close(startedCh)
		defer func() {
			logger.Debugf("Receiver's '%s' port observer stopped", input.Name)
			if err := inMsgSubs.Unsubscribe(); err != nil {
				panic(err)
			}
			close(inMsgCh)
			wg.Done()
		}()

		for {
			select {
			case <-doneCh:
				logger.Debugf("Receiver's '%s' port observer shut down", input.Name)
				return

			case inputMsg := <-inMsgCh:
				logger.Debugf("Receiver's '%s' port observer received message", input.Name)
				newInput := io.NewInput(input.Name, input.Type, input.Representation, input.Channel, input.DefaultMessage)
				newInput.Message = msgs.GetDefaultMessageByType(input.Type)
				if err := newInput.Message.Decode(input.Representation, inputMsg); err != nil {
					panic(err)
				}
				inputsMuxCh <- newInput
				logger.Debugf("Receiver's '%s' port observer sent message to inputMuxCh channel", input.Name)
			}
		}
	}()
	return startedCh
}
