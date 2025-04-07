package services

import (
	"testing"
	"time"

	"go-authentication/internal/snmp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNatsService struct {
	mock.Mock
}

func (m *MockNatsService) Publish(subject string, data []byte) error {
	args := m.Called(subject, data)
	return args.Error(0)
}

func (m *MockNatsService) Close() {
	m.Called()
}

func TestSNMPService(t *testing.T) {
	config := &snmp.SNMPConfig{
		Port:      161,
		Community: "public",
		Interval:  100 * time.Millisecond,
		Devices: []snmp.MockDevice{
			{
				ID:          "device1",
				IP:          "192.168.1.1",
				DataPattern: "random",
				OIDs:        []string{"1.3.6.1.2.1.1.1.0"},
			},
		},
	}

	simulator := snmp.NewSNMPSimulator(config)
	mockNats := &MockNatsService{}

	// Set up expectations
	mockNats.On("Publish", "snmp.data", mock.Anything).Return(nil)
	mockNats.On("Close").Return()

	service := NewSNMPService(simulator, mockNats)
	assert.NotNil(t, service)

	// Start publishing
	service.StartPublishing()

	// Wait for a publish cycle
	time.Sleep(150 * time.Millisecond)

	// Verify that Publish was called
	mockNats.AssertExpectations(t)
}
