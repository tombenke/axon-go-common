// Package outputs provides the functions to forward the results of the processing
package outputs

import (
	"github.com/sirupsen/logrus"
	"github.com/tombenke/axon-go-common/io"
	"github.com/tombenke/axon-go-common/messenger"
	"github.com/tombenke/axon-go-common/msgs"
	"github.com/tombenke/axon-go-common/msgs/orchestra"
	"sync"
)

// SyncSender receives outputs from the processor function via the `outputsCh` that it sends to
// the corresponding topics identified by the port.
// The outputs structures hold every details about the ports, the message itself, and the subject to send.
// This function runs as a standalone process, so it should be started as a go function.
func SyncSender(actorName string, outputsCh chan io.Outputs, doneCh chan interface{}, wg *sync.WaitGroup, m messenger.Messenger, logger *logrus.Logger) (chan interface{}, chan interface{}) {
	var outputs io.Outputs
	senderStoppedCh := make(chan interface{})
	startedCh := make(chan interface{})

	wg.Add(1)
	go func() {
		sendResultsCh := make(chan []byte)
		sendResultsSubs := m.ChanSubscribe("send-results", sendResultsCh)
		logger.Debugf("Sender started in sync mode.")
		close(startedCh)

		defer func() {
			if err := sendResultsSubs.Unsubscribe(); err != nil {
				panic(err)
			}
			close(sendResultsCh)
			logger.Debugf("Sender stopped")
			close(senderStoppedCh)
			wg.Done()
		}()

		for {
			select {
			case <-doneCh:
				logger.Debugf("Sender shuts down.")
				return

			case outputs = <-outputsCh:
				logger.Debugf("Sender received outputs")
				// In sync mode notifies the orchestrator about that it is ready to send
				sendProcessingCompleted(actorName, m)

			case <-sendResultsCh:
				logger.Debugf("Sender received orchestrator trigger to send outputs")
				syncSendOutputs(actorName, outputs, m)
			}
		}
	}()

	return startedCh, senderStoppedCh
}

// sendProcessingCompleted sends a message to the orchestrator about that
// the agent completed the processing and it is ready to send outputs.
func sendProcessingCompleted(actorName string, m messenger.Messenger) {
	logger.Debugf("Sender sends 'processing-completed' notification to orchestrator\n")
	processingCompletedMsg := orchestra.NewProcessingCompletedMessage(actorName)
	if err := m.Publish("processing-completed", processingCompletedMsg.Encode(msgs.JSONRepresentation)); err != nil {
		panic(err)
	}
}

func syncSendOutputs(actorName string, outputs io.Outputs, m messenger.Messenger) {
	for o := range outputs {
		channel := outputs[o].Channel
		representation := outputs[o].Representation
		message := outputs[o].Message
		messageType := outputs[o].Type
		logger.Debugf("Sender sends '%v' type message of '%s' output port to '%s' channel in '%s' format\n", messageType, o, channel, representation)
		if err := m.Publish(channel, message.Encode(representation)); err != nil {
			panic(err)
		}
	}

	logger.Debugf("Sender sends 'sending-completed' notification to orchestrator\n")
	sendingCompletedMsg := orchestra.NewSendingCompletedMessage(actorName)
	if err := m.Publish("sending-completed", sendingCompletedMsg.Encode(msgs.JSONRepresentation)); err != nil {
		panic(err)
	}
}
