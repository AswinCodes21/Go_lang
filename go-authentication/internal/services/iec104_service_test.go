package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIEC104Service(t *testing.T) {
	// Create IEC104 service without NATS
	service := &IEC104Service{
		port:    502,
		timeout: 10,
		k:       1,
		w:       1,
		devices: make([]*IEC104Device, 0),
	}

	// Test GetDevices
	t.Run("GetDevices", func(t *testing.T) {
		devices := service.GetDevices()
		assert.NotNil(t, devices)
		assert.Equal(t, 0, len(devices))
	})

	// Test GetLatestData
	t.Run("GetLatestData", func(t *testing.T) {
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
}
