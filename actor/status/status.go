// Package status provides the functions to communicate the status of the actor with the orchestrator application
package status

import (
	"github.com/sirupsen/logrus"
	"github.com/tombenke/axon-go-common/config"
	"github.com/tombenke/axon-go-common/messenger"
	"github.com/tombenke/axon-go-common/msgs"
	"github.com/tombenke/axon-go-common/msgs/orchestra"
	"sync"
)

// Status receives status request messages from the orchestrator application,
// sends responses to these requests, forwarding the actual status of the actor.
// This function runs as a standalone process, so it should be started as a go function.
func Status(nodeConfig config.Node, doneCh chan bool, wg *sync.WaitGroup, m messenger.Messenger, logger *logrus.Logger) (chan bool, chan bool) {
	statusRequestCh := make(chan []byte)
	statusRequestSubs := m.ChanSubscribe(nodeConfig.Orchestration.Channels.StatusRequest, statusRequestCh)
	statusStoppedCh := make(chan bool)
	statusStartedCh := make(chan bool)

	wg.Add(1)
	go func() {
		close(statusStartedCh)
		defer func() {
			if err := statusRequestSubs.Unsubscribe(); err != nil {
				panic(err)
			}
			close(statusRequestCh)
			wg.Done()

			logger.Debugf("Status stopped.")
			close(statusStoppedCh)
		}()

		for {
			select {
			case <-doneCh:
				logger.Debugf("Status shuts down.")
				return

			case <-statusRequestCh:
				logger.Debugf("Status received status-request message")
				logger.Debugf("Status sends status-report message")
				statusReportMsg := makeStatusReportMsg(nodeConfig)
				if err := m.Publish(nodeConfig.Orchestration.Channels.StatusReport, statusReportMsg.Encode(msgs.JSONRepresentation)); err != nil {
					panic(err)
				}
				// TODO: Make orchestra message representations configurable
			}
		}
	}()
	logger.Debugf("Status started")
	return statusStartedCh, statusStoppedCh
}

func makeStatusReportMsg(nodeConfig config.Node) msgs.Message {
	srBody := orchestra.StatusReportBody{

		Name:            nodeConfig.Name,
		Type:            nodeConfig.Type,
		Ports:           makePorts(nodeConfig.Ports),
		Synchronization: nodeConfig.Orchestration.Synchronization,
		SpecsURL:        nodeConfig.SpecsURL,
	}

	return orchestra.NewStatusReportMessage(srBody)
}

func makePorts(portsConfig config.Ports) orchestra.Ports {
	ports := orchestra.Ports{
		Inputs:  make([]orchestra.Port, 0),
		Outputs: make([]orchestra.Port, 0),
	}

	for _, in := range portsConfig.Inputs {
		ports.Inputs = append(ports.Inputs, orchestra.Port{
			Name:           in.Name,
			Type:           in.Type,
			Representation: in.Representation,
			Channel: orchestra.Channel{
				Name: in.Channel,
				Type: "TOPIC",
			},
		})
	}

	for _, out := range portsConfig.Outputs {
		ports.Outputs = append(ports.Outputs, orchestra.Port{
			Name:           out.Name,
			Type:           out.Type,
			Representation: out.Representation,
			Channel: orchestra.Channel{
				Name: out.Channel,
				Type: "TOPIC",
			},
		})
	}

	return ports
}
