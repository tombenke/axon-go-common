// Package inputs provides the functions to receive input messages, collect them and forward to the processor
package inputs

import (
	"github.com/sirupsen/logrus"
	"github.com/tombenke/axon-go-common/config"
	"github.com/tombenke/axon-go-common/io"
	"github.com/tombenke/axon-go-common/messenger"
	"sync"
)

// AsyncReceiver receives inputs from the connecting actors processor function via the `outputsCh`
// that it sends to the processor for further processing.
// The inputs structures hold every details about the ports, the message itself,
// and the subject to receive from.
// This function starts the receiver routine as a standalone process,
// and returns a channel that the process uses to forward the incoming inputs.
func AsyncReceiver(inputsCfg config.Inputs, resetCh chan interface{}, doneCh chan interface{}, appWg *sync.WaitGroup, m messenger.Messenger, logger *logrus.Logger) (chan interface{}, chan *io.Inputs, chan interface{}) {
	receiverStoppedCh := make(chan interface{})
	startedCh := make(chan interface{})

	// Setup communication channel with the processor
	inputsCh := make(chan *io.Inputs)

	appWg.Add(1)
	go func() {
		logger.Debugf("Receiver started in async mode.")
		close(startedCh)
		defer logger.Debugf("Receiver stopped.")
		defer close(inputsCh)
		defer close(receiverStoppedCh)
		defer appWg.Done()

		// Create wait-group for the channel observer sub-processes
		obsWg := sync.WaitGroup{}
		obsDoneCh := make(chan interface{})

		// Create Input ports, and initialize with default messages
		inputs := asyncSetupInputPorts(inputsCfg, logger)

		// Creates an inputs multiplexer channel for observers to send their inputs via one channel
		inputsMuxCh := make(chan io.Input)
		defer close(inputsMuxCh)

		// Starts the input port observers
		startInPortsObservers(inputs, inputsMuxCh, obsDoneCh, &obsWg, m, logger)

		for {
			select {
			case <-doneCh:
				logger.Debugf("Receiver shuts down.")
				close(obsDoneCh)
				logger.Debugf("Receiver closed the 'obsDoneCh'.")
				logger.Debugf("Receiver starts waiting for observers to stop")
				obsWg.Wait()
				logger.Debugf("Receiver's observers stopped")
				return

			case <-resetCh:
				logger.Debugf("Receiver got RESET signal")
				inputsCh <- inputs
				logger.Debugf("Receiver sent 'inputs' to 'inputsCh'")

			case input := <-inputsMuxCh:
				logger.Debugf("Receiver got message to '%s' port", input.Name)
				(*inputs).SetMessage(input.Name, input.Message)
				// Immediately forward to the processor if not in synchronized mode
				inputsCh <- inputs
				logger.Debugf("Receiver sent 'inputs' to 'inputsCh'")
			}
		}
	}()

	return startedCh, inputsCh, receiverStoppedCh
}

// asyncSetupInputPorts creates inputs ports, and initilizes them with their default messages
func asyncSetupInputPorts(inputsCfg config.Inputs, logger *logrus.Logger) *io.Inputs {

	logger.Debugf("Receiver sets up input ports")

	// Create input ports
	inputs := io.NewInputs(inputsCfg)

	// Set every input ports' message to its default
	for p := range inputs.Map {
		defaultMessage := (*inputs).Map[p].DefaultMessage
		inputs.SetMessage(p, defaultMessage)
	}

	return inputs
}
