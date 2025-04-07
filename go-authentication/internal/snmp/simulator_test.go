package snmp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSNMPSimulator(t *testing.T) {
	config := &SNMPConfig{
		Port:      161,
		Community: "public",
		Interval:  30 * time.Second,
		Devices: []MockDevice{
			{
				ID:          "device1",
				IP:          "192.168.1.1",
				DataPattern: "random",
				OIDs:        []string{"1.3.6.1.2.1.1.1.0", "1.3.6.1.2.1.1.3.0"},
			},
		},
	}

	simulator := NewSNMPSimulator(config)

	assert.NotNil(t, simulator)
	assert.Equal(t, config, simulator.Config)
	assert.Len(t, simulator.Devices, 1)
	assert.Equal(t, "device1", simulator.Devices[0].ID)
}

func TestGenerateMockData(t *testing.T) {
	config := &SNMPConfig{
		Port:      161,
		Community: "public",
		Interval:  30 * time.Second,
		Devices: []MockDevice{
			{
				ID:          "device1",
				IP:          "192.168.1.1",
				DataPattern: "random",
				OIDs:        []string{"1.3.6.1.2.1.1.1.0", "1.3.6.1.2.1.1.3.0"},
			},
		},
	}

	simulator := NewSNMPSimulator(config)
	data := simulator.GenerateMockData(simulator.Devices[0])

	assert.NotNil(t, data)
	assert.Len(t, data, 2)
	assert.Equal(t, "device1", data[0].DeviceID)
	assert.NotNil(t, data[0].Value)
	assert.NotNil(t, data[0].Timestamp)
}

func TestStart(t *testing.T) {
	config := &SNMPConfig{
		Port:      161,
		Community: "public",
		Interval:  30 * time.Second,
		Devices:   []MockDevice{},
	}

	simulator := NewSNMPSimulator(config)
	err := simulator.Start()

	assert.NoError(t, err)
} 