package services

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

type IEC104Service struct {
	port        int
	timeout     int
	k           int
	w           int
	natsService NatsServiceInterface
	devices     []*IEC104Device
	mu          sync.RWMutex
}

type IEC104Device struct {
	ID         string
	Name       string
	Address    string
	LastSeen   time.Time
	Status     string
	DataPoints []DataPoint
}

type DataPoint struct {
	Address int
	Value   float64
	Quality int
	Time    time.Time
}

type IEC104Message struct {
	DeviceID  string      `json:"device_id"`
	Timestamp time.Time   `json:"timestamp"`
	Data      []DataPoint `json:"data"`
}

func NewIEC104Service(port, timeout, k, w int, natsService NatsServiceInterface) *IEC104Service {
	service := &IEC104Service{
		port:        port,
		timeout:     timeout,
		k:           k,
		w:           w,
		natsService: natsService,
		devices:     make([]*IEC104Device, 0),
	}

	// Add sample devices with more data points
	service.mu.Lock()
	defer service.mu.Unlock()

	service.devices = append(service.devices, &IEC104Device{
		ID:       "1",
		Name:     "RTU-1",
		Address:  "192.168.1.100",
		LastSeen: time.Now(),
		Status:   "disconnected",
		DataPoints: []DataPoint{
			{
				Address: 1,
				Value:   42.5,
				Quality: 0,
				Time:    time.Now(),
			},
			{
				Address: 2,
				Value:   38.2,
				Quality: 0,
				Time:    time.Now(),
			},
			{
				Address: 3,
				Value:   15.7,
				Quality: 0,
				Time:    time.Now(),
			},
		},
	})

	service.devices = append(service.devices, &IEC104Device{
		ID:       "2",
		Name:     "RTU-2",
		Address:  "192.168.1.101",
		LastSeen: time.Now(),
		Status:   "disconnected",
		DataPoints: []DataPoint{
			{
				Address: 1,
				Value:   25.3,
				Quality: 0,
				Time:    time.Now(),
			},
			{
				Address: 2,
				Value:   19.8,
				Quality: 0,
				Time:    time.Now(),
			},
			{
				Address: 3,
				Value:   45.6,
				Quality: 0,
				Time:    time.Now(),
			},
		},
	})

	log.Printf("IEC104Service initialized with %d devices", len(service.devices))
	return service
}

func (s *IEC104Service) Start() error {
	// Initialize simulator connection
	go s.simulateData()
	return nil
}

func (s *IEC104Service) simulateData() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		for _, device := range s.devices {
			// Simulate data changes
			for i := range device.DataPoints {
				// Random value between 0 and 100
				device.DataPoints[i].Value = rand.Float64() * 100
				device.DataPoints[i].Time = time.Now()
				device.DataPoints[i].Quality = 0
			}

			// Create IEC104 message
			message := IEC104Message{
				DeviceID:  device.ID,
				Timestamp: time.Now(),
				Data:      device.DataPoints,
			}

			// Marshal to JSON
			msgBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling data for device %s: %v", device.ID, err)
				continue
			}

			// Publish data through NATS
			subject := fmt.Sprintf("iec104.device.%s.data", device.ID)
			if err := s.natsService.Publish(subject, msgBytes); err != nil {
				log.Printf("Failed to publish data for device %s: %v", device.ID, err)
			} else {
				log.Printf("Published data for device %s", device.ID)
			}
		}
		s.mu.Unlock()
	}
}

func (s *IEC104Service) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Set connection timeout
	conn.SetDeadline(time.Now().Add(time.Duration(s.timeout) * time.Second))

	// TODO: Implement IEC 104 protocol handling
	// 1. Handle STARTDT
	// 2. Process ASDU messages
	// 3. Handle TESTFR
	// 4. Manage sequence numbers
}

func (s *IEC104Service) publishData(device *IEC104Device) {
	message := IEC104Message{
		DeviceID:  device.ID,
		Timestamp: time.Now(),
		Data:      device.DataPoints,
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling IEC 104 data: %v", err)
		return
	}

	err = s.natsService.Publish("iec104.data", msgBytes)
	if err != nil {
		log.Printf("Error publishing IEC 104 data: %v", err)
	}
}

func (s *IEC104Service) GetStatus() (bool, int) {
	return true, len(s.devices)
}

func (s *IEC104Service) GetDevices() []*IEC104Device {
	s.mu.RLock()
	defer s.mu.RUnlock()

	log.Printf("Getting devices, count: %d", len(s.devices))
	for _, device := range s.devices {
		log.Printf("Device found: ID=%s, Name=%s, Status=%s", device.ID, device.Name, device.Status)
	}
	return s.devices
}

func (s *IEC104Service) GetLatestData(deviceID string) []DataPoint {
	for _, device := range s.devices {
		if device.ID == deviceID {
			return device.DataPoints
		}
	}
	return nil
}

func (s *IEC104Service) GetDeviceStatus(deviceID string) string {
	for _, device := range s.devices {
		if device.ID == deviceID {
			return device.Status
		}
	}
	return "disconnected"
}

func (s *IEC104Service) ConnectToDevice(deviceID string) {
	for _, device := range s.devices {
		if device.ID == deviceID {
			device.Status = "connected"
			device.LastSeen = time.Now()
			return
		}
	}
}

func (s *IEC104Service) DisconnectFromDevice(deviceID string) {
	for _, device := range s.devices {
		if device.ID == deviceID {
			device.Status = "disconnected"
			return
		}
	}
}
