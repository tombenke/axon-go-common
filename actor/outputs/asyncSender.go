// Package outputs provides the functions to forward the results of the processing
package outputs

import (
	"github.com/sirupsen/logrus"
	"github.com/tombenke/axon-go-common/io"
	"github.com/tombenke/axon-go-common/messenger"
	"sync"
)

// AsyncSender receives outputs from the processor function via the `outputsCh` that it sends to
// the corresponding topics identified by the port.
// The outputs structures hold every details about the ports, the message itself, and the subject to send.
// This function runs as a standalone process, so it should be started as a go function.
func AsyncSender(actorName string, outputsCh chan io.Outputs, doneCh chan interface{}, wg *sync.WaitGroup, m messenger.Messenger, logger *logrus.Logger) (chan interface{}, chan interface{}) {
	var outputs io.Outputs
	senderStoppedCh := make(chan interface{})
	startedCh := make(chan interface{})

	wg.Add(1)
	go func() {
		logger.Debugf("Sender started in async mode.")
		close(startedCh)
		defer logger.Debugf("Sender stopped")
		defer close(senderStoppedCh)
		defer wg.Done()

		for {
			select {
			case <-doneCh:
				logger.Debugf("Sender shuts down.")
				return

			case outputs = <-outputsCh:
				logger.Debugf("Sender received outputs")
				// In async mode it immediately sends the outputs whet it gets them
				asyncSendOutputs(actorName, outputs, m, logger)
			}
		}
	}()

	return startedCh, senderStoppedCh
}

func asyncSendOutputs(actorName string, outputs io.Outputs, m messenger.Messenger, logger *logrus.Logger) {
	for o := range outputs {
		message := outputs[o].Message
		channel := outputs[o].Channel
		representation := outputs[o].Representation
		messageType := outputs[o].Type
		if message != nil {
			logger.Debugf("Sender sends '%v' type message of '%s' output port to '%s' channel in '%s' format", messageType, o, channel, representation)
			if err := m.Publish(channel, message.Encode(representation)); err != nil {
				panic(err)
			}
		} else {
			logger.Errorf("Sender wants to send '%v' type message of '%s' output port to '%s' channel in '%s' format but message is nil", messageType, o, channel, representation)
		}
	}
}
