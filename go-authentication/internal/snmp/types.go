package snmp

import "time"

// SNMPData represents a single SNMP data point
type SNMPData struct {
	OID       string      `json:"oid"`
	Value     interface{} `json:"value"`
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	DeviceID  string      `json:"device_id"`
}

// MockDevice represents a simulated SNMP device
type MockDevice struct {
	ID          string   `json:"id"`
	IP          string   `json:"ip"`
	Community   string   `json:"community"`
	OIDs        []string `json:"oids"`
	DataPattern string   `json:"data_pattern"`
}

// SNMPConfig holds SNMP configuration
type SNMPConfig struct {
	Community string        `json:"community"`
	Port      int           `json:"port"`
	Interval  time.Duration `json:"interval"`
	Devices   []MockDevice  `json:"devices"`
}
