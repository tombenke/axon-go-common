package config

import (
	////"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tombenke/axon-go-common/messenger"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

type AppConfigFromFile struct {
	Node           Node   `yaml:"node"`
	ExtDescription string `yaml:"extDescription"`
}

// YAML converts the content of the AppConfigFromFile structure to YAML format
func (c AppConfigFromFile) YAML() ([]byte, error) {
	return yaml.Marshal(&c)
}

func TestReadConfigFile(t *testing.T) {

	expectedAppConfigFromFile := AppConfigFromFile{
		Node: Node{
			Messenger: messenger.Config{
				Urls:      defaultMessagingURL,
				UserCreds: defaultMessagingUserCreds,
			},
			Name:           "well-water-upper-level-sensor-simulator",
			Type:           "",
			ConfigFileName: "config.yml",
			LogLevel:       "debug",
			LogFormat:      "text",
			Ports: Ports{
				Inputs: Inputs{
					In{
						IO: IO{
							Name:           "reference-water-level",
							Type:           "base/Float64",
							Representation: "application/json",
							Channel:        "",
						},
						Default: "{ \"Body\": { \"Data\": 0.75 } }",
					},
					In{
						IO: IO{
							Name:           "water-level",
							Type:           "base/Float64",
							Representation: "application/json",
							Channel:        "well-water-level",
						},
						Default: "{ \"Body\": { \"Data\": 0.0 } }",
					},
				},
				Outputs: Outputs{
					Out{
						IO: IO{
							Name:           "water-level-state",
							Type:           "base/Bool",
							Representation: "application/json",
							Channel:        "well-water-upper-level-state",
						},
					},
				},
				Configure: Configure{
					Extend: false,
					Modify: true,
				},
			},
			Orchestration: Orchestration{
				Presence:        true,
				Synchronization: true,
				Channels: Channels{
					StatusRequest:       "status-request",
					StatusReport:        "status-report",
					SendResults:         "send-results",
					SendingCompleted:    "sending-completed",
					ReceiveAndProcess:   "receive-and-process",
					ProcessingCompleted: "processing-completed",
				},
			},
		},
		ExtDescription: "This is an extensional property",
	}
	//emptyAppConfig := AppConfig{}
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path := cwd + "/test-config.yml"
	appConfigFromFile, err := ReadAppConfigFromFile(path)
	assert.Nil(t, err)
	////appConfigFromFileYAML, err := appConfigFromFile.YAML()
	////fmt.Printf("appConfigFromFile:\n%s\n", appConfigFromFileYAML)
	assert.Nil(t, err)

	assert.Equal(t, expectedAppConfigFromFile, appConfigFromFile)
}

// ReadConfigFromFile Reads the the config parameters from a file
func ReadAppConfigFromFile(path string) (AppConfigFromFile, error) {
	c := AppConfigFromFile{}
	var err error
	content, err := LoadFile(path)
	if err == nil {
		err = yaml.Unmarshal([]byte(content), &c)
	}
	return c, err
}
