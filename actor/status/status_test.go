package status

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/tombenke/axon-go-common/config"
	"github.com/tombenke/axon-go-common/messenger"
	messengerImpl "github.com/tombenke/axon-go-common/messenger/nats"
	"github.com/tombenke/axon-go-common/msgs"
	"github.com/tombenke/axon-go-common/msgs/orchestra"
	at "github.com/tombenke/axon-go-common/testing"
	"sync"
	"testing"
)

const (
	checkSendStatusRequest    = "orchestrator sent status request"
	checkStatusReportReceived = "orchestrator received status report"
)

var checklist = []string{
	checkSendStatusRequest,
	checkStatusReportReceived,
}

var logger = logrus.New()

var messengerCfg = messenger.Config{
	Urls:       "localhost:4222",
	UserCreds:  "",
	ClientName: "status-test-client",
	ClusterID:  "test-cluster",
	ClientID:   "status-test-client",
	Logger:     logger,
}

func createTestNode() config.Node {
	node := config.NewNode("test-node", "test-node-type", true, true, true, false)

	node.AddInputPort("input", "base/Any", "application/json", "axon-function-js.input", "")
	node.AddOutputPort("output", "base/Any", "application/json", "axon-function-js.output")

	node.AddSpecsURL("https://github.com/tombenke/axon-go")

	return node
}

func configToStatusReportBody(node config.Node) orchestra.StatusReportBody {
	srBody := orchestra.StatusReportBody{
		Name:            node.Name,
		Type:            node.Type,
		Synchronization: node.Orchestration.Synchronization,
		SpecsURL:        node.SpecsURL,
		Ports: orchestra.Ports{
			Inputs:  make([]orchestra.Port, 0),
			Outputs: make([]orchestra.Port, 0),
		},
	}

	for _, in := range node.Ports.Inputs {
		srBody.Ports.Inputs = append(srBody.Ports.Inputs, orchestra.Port{
			Name:           in.Name,
			Type:           in.Type,
			Representation: in.Representation,
			Channel: orchestra.Channel{
				Name: in.Channel,
				Type: "TOPIC",
			},
		})
	}

	for _, out := range node.Ports.Outputs {
		srBody.Ports.Outputs = append(srBody.Ports.Outputs, orchestra.Port{
			Name:           out.Name,
			Type:           out.Type,
			Representation: out.Representation,
			Channel: orchestra.Channel{
				Name: out.Channel,
				Type: "TOPIC",
			},
		})
	}

	return srBody
}

func TestStatus(t *testing.T) {
	testNode := createTestNode()
	expectedStatusReportBody := configToStatusReportBody(testNode)

	// Connect to messaging
	m := messengerImpl.NewMessenger(messengerCfg)
	defer m.Close()

	// Use a WaitGroup to wait for the processes of the testbed to complete their mission
	wg := sync.WaitGroup{}

	// Create a trigger channel to start the test
	triggerCh := make(chan interface{})

	// Start the processes of the test-bed
	doneChkCh := make(chan interface{})
	reportCh, testCompletedCh, checklistStoppedCh := at.ChecklistProcess(checklist, doneChkCh, &wg, logger)

	doneOrcCh := make(chan interface{})
	orcStoppedCh := startMockOrchestrator(t, reportCh, triggerCh, doneOrcCh, &wg, logger, m, expectedStatusReportBody)

	// Start the status process
	doneStatusCh := make(chan interface{})
	statusStartedCh, statusStoppedCh := Status(testNode, doneStatusCh, &wg, m, logger)

	// Wait until all components have been successfully started
	<-statusStartedCh

	// Start testing
	triggerCh <- true

	// Wait until test is completed, then stop the processes
	logger.Infof("Wait until test is completed")
	<-testCompletedCh

	logger.Infof("Stops Orchestrator")
	close(doneOrcCh)
	logger.Infof("Wait Orchestrator to stop")
	<-orcStoppedCh
	logger.Infof("Orchestrator stopped")

	logger.Infof("Stops Status")
	close(doneStatusCh)
	logger.Infof("Wait Status to stop")
	<-statusStoppedCh
	logger.Infof("Status stopped")

	logger.Infof("Stops Checklist")
	close(doneChkCh)
	logger.Infof("Wait Checklist to stop")
	<-checklistStoppedCh
	logger.Infof("Checklist stopped")

	// Wait for the message to come in
	wg.Wait()
}

// startMockOrchestrator starts a standalone process that emulates
// the behaviour of an external orchestrator application.
// Orchestrator waits for an incoming trigger to start the test process via sending a status request,
// then waits for receiving the status response.
// The Mock Orchestrator reports every relevant event to the Checklist process.
// Mock Orchestrator will shut down if it receives a message via the `doneCh` channel.
func startMockOrchestrator(t *testing.T, reportCh chan string, triggerCh chan interface{}, doneCh chan interface{}, wg *sync.WaitGroup, logger *logrus.Logger, m messenger.Messenger, expectedStatusReportBody orchestra.StatusReportBody) chan interface{} {
	statusReportCh := make(chan []byte)
	statusReportSubs := m.ChanSubscribe("status-report", statusReportCh)
	orcStoppedCh := make(chan interface{})

	wg.Add(1)
	go func() {
		defer func() {
			logger.Infof("MockOrchestrator stopped.")
			if err := statusReportSubs.Unsubscribe(); err != nil {
				panic(err)
			}
			close(statusReportCh)
			close(orcStoppedCh)
			wg.Done()
		}()

		for {
			select {
			case <-doneCh:
				logger.Infof("MockOrchestrator shuts down.")
				return

			case <-triggerCh:
				logger.Infof("MockOrchestrator received 'start-trigger'.")
				logger.Infof("MockOrchestrator sends 'status-request' message.")
				statusRequestMsg := orchestra.NewStatusRequestMessage()
				if err := m.Publish("status-request", statusRequestMsg.Encode(msgs.JSONRepresentation)); err != nil {
					panic(err)
				}
				// TODO: Make orchestra message representations and channel names configurable
				reportCh <- checkSendStatusRequest

			case statusReportMsgBytes := <-statusReportCh:
				logger.Infof("MockOrchestrator received 'status-report' message.")
				var statusReportMsg orchestra.StatusReport
				err := statusReportMsg.Decode(msgs.JSONRepresentation, statusReportMsgBytes)
				assert.Nil(t, err)
				assert.Equal(t, expectedStatusReportBody, statusReportMsg.Body)
				reportCh <- checkStatusReportReceived
			}
		}
	}()
	logger.Infof("Mock Orchestrator started.")

	return orcStoppedCh
}
