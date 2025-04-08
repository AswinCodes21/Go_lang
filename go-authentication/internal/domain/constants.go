package domain

const (
	// Default port for IEC104 protocol
	DefaultPort = 2404

	// Device status constants
	StatusConnected    = "connected"
	StatusDisconnected = "disconnected"
	StatusError        = "error"

	// Error messages
	ErrDeviceNotFound = "device not found"
	ErrInvalidPort    = "invalid port number"
	ErrConnection     = "connection error"
	ErrTimeout        = "connection timeout"
	ErrInvalidData    = "invalid data format"

	// NATS subjects
	SubjectDeviceStatus = "iec104.device.status"
	SubjectDeviceData   = "iec104.device.data"

	// API routes
	RouteDevices          = "/iec104/devices"
	RouteDeviceStatus     = "/iec104/device/:device_id/status"
	RouteDeviceData       = "/iec104/device/:device_id/data"
	RouteDeviceConnect    = "/iec104/device/:device_id/connect"
	RouteDeviceDisconnect = "/iec104/device/:device_id/disconnect"
)
