package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNatsService is a mock implementation of NatsServiceInterface
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

func TestIEC104Service(t *testing.T) {
	// Create mock NATS service
	mockNats := new(MockNatsService)
	
	// Create IEC104 service with default parameters
	service := NewIEC104Service(2404, 30, 12, 8, mockNats)

	// Test GetDevices
	t.Run("GetDevices", func(t *testing.T) {
		devices := service.GetDevices()
		assert.NotNil(t, devices)
		assert.Equal(t, 0, len(devices)) // Initially empty
	})

	// Test GetLatestData
	t.Run("GetLatestData", func(t *testing.T) {
		// Test with invalid device ID
		data := service.GetLatestData("invalid-id")
		assert.Nil(t, data)
	})

	// Test GetStatus
	t.Run("GetStatus", func(t *testing.T) {
		status, count := service.GetStatus()
		assert.True(t, status)
		assert.Equal(t, 0, count)
	})

	// Test Start
	t.Run("Start", func(t *testing.T) {
		err := service.Start()
		assert.NoError(t, err)
	})

	// Test publishData
	t.Run("PublishData", func(t *testing.T) {
		device := &IEC104Device{
			ID:        "test-device",
			Name:      "Test Device",
			Address:   "127.0.0.1",
			LastSeen:  time.Now(),
			Status:    "connected",
			DataPoints: []DataPoint{
				{
					Address: 1,
					Value:   42.0,
					Quality: 0,
					Time:    time.Now(),
				},
			},
		}

		// Set up mock expectations
		mockNats.On("Publish", "iec104.data", mock.Anything).Return(nil)

		// Call publishData
		service.publishData(device)

		// Verify that Publish was called
		mockNats.AssertExpectations(t)
	})
} 