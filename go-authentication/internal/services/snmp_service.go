package services

import (
	"encoding/json"
	"log"
	"time"

	"go-authentication/internal/snmp"
)

type SNMPService struct {
	simulator   *snmp.SNMPSimulator
	natsService NatsServiceInterface
}

type SNMPMessage struct {
	DeviceID  string          `json:"device_id"`
	Timestamp time.Time       `json:"timestamp"`
	Data      []snmp.SNMPData `json:"data"`
}

func NewSNMPService(simulator *snmp.SNMPSimulator, natsService NatsServiceInterface) *SNMPService {
	return &SNMPService{
		simulator:   simulator,
		natsService: natsService,
	}
}

func (s *SNMPService) StartPublishing() {
	ticker := time.NewTicker(s.simulator.Config.Interval)
	go func() {
		for range ticker.C {
			s.publishSNMPData()
		}
	}()
}

func (s *SNMPService) publishSNMPData() {
	for _, device := range s.simulator.Devices {
		data := s.simulator.GenerateMockData(device)
		message := SNMPMessage{
			DeviceID:  device.ID,
			Timestamp: time.Now(),
			Data:      data,
		}

		msgBytes, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshaling SNMP data: %v", err)
			continue
		}

		err = s.natsService.Publish("snmp.data", msgBytes)
		if err != nil {
			log.Printf("Error publishing SNMP data: %v", err)
		}
	}
}

// GetStatus returns the current status of the SNMP service
func (s *SNMPService) GetStatus() (bool, int) {
	if s.simulator == nil {
		return false, 0
	}
	return true, len(s.simulator.Devices)
}

// GetDevices returns the list of simulated devices
func (s *SNMPService) GetDevices() []*snmp.MockDevice {
	if s.simulator == nil {
		return nil
	}
	return s.simulator.Devices
}

// GetLatestData returns the latest data for a device
func (s *SNMPService) GetLatestData(deviceID string) []snmp.SNMPData {
	if s.simulator == nil {
		return nil
	}

	// Find the device
	var device *snmp.MockDevice
	for _, d := range s.simulator.Devices {
		if d.ID == deviceID {
			device = d
			break
		}
	}

	if device == nil {
		return nil
	}

	// Generate latest data
	return s.simulator.GenerateMockData(device)
}
