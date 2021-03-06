package inputs

import (
	messengerImpl "github.com/tombenke/axon-go-common/messenger/nats"
	at "github.com/tombenke/axon-go-common/testing"
	"sync"
	"testing"
	"time"
)

// TestAsyncReceiverStartStop sets up the input ports, then stops.
func TestAsyncReceiverStartStop(t *testing.T) {
	// Connect to messaging
	m := messengerImpl.NewMessenger(messengerCfg)
	defer m.Close()

	// Use a WaitGroup to wait for the processes of the testbed to complete their mission
	wg := sync.WaitGroup{}

	// Create a channel for the RESET
	resetCh := make(chan interface{})

	// Create a channel to shut down the processes if needed
	doneCh := make(chan interface{})

	// Start the receiver process
	startedCh, _, _ := AsyncReceiver(asyncInputsCfg, resetCh, doneCh, &wg, m, logger)
	<-startedCh

	// Wait until test is completed, then stop the processes
	close(doneCh)

	// Wait for the message to come in
	wg.Wait()
	close(resetCh)
}

// TestReceiveInputs sets up the input ports, and gets inputs to each ports, then a receive-and-process message,
// It uses the incoming messages that it sends as the result inputs to the processor.
func TestAsyncReceiverInputs(t *testing.T) {
	// Connect to messaging
	m := messengerImpl.NewMessenger(messengerCfg)
	defer m.Close()

	// Use a WaitGroup to wait for the processes of the testbed to complete their mission
	wg := sync.WaitGroup{}

	// Start the processes of the test-bed
	doneCheckCh := make(chan interface{})
	reportCh, testCompletedCh, chkStoppedCh := at.ChecklistProcess(checklistAsync, doneCheckCh, &wg, logger)

	// Create a channel for the RESET
	resetCh := make(chan interface{})

	// Start the receiver process
	doneRcvCh := make(chan interface{})
	startedCh, inputsCh, rcvStoppedCh := AsyncReceiver(asyncInputsCfg, resetCh, doneRcvCh, &wg, m, logger)
	<-startedCh

	doneProcCh := make(chan interface{})
	procStoppedCh := startMockProcessor(inputsCh, reportCh, doneProcCh, &wg, logger)

	// Give chance for observers to start before send messages through external messaging mw.
	time.Sleep(100 * time.Millisecond)

	// Start testing
	sendInputMessages(asyncInputsCfg, asyncInputs, reportCh, m, logger)

	// Wait until test is completed, then stop the processes
	<-testCompletedCh

	logger.Infof("Stops Mock Processor")
	close(doneProcCh)
	logger.Infof("Wait Mock Processor to stop")
	<-procStoppedCh
	logger.Infof("Mock Processor stopped")

	logger.Infof("Stops Receiver")
	close(doneRcvCh)
	logger.Infof("Wait Receiver to stop")
	<-rcvStoppedCh
	close(resetCh)
	logger.Infof("Receiver stopped")

	logger.Infof("Stops Checklist")
	close(doneCheckCh)
	logger.Infof("Wait Checklist to stop")
	<-chkStoppedCh
	logger.Infof("Checklist stopped")

	// Wait for the message to come in
	wg.Wait()
}
