package constants

const (
	// IEC104 Protocol Constants
	IEC104DefaultPort    = 2404
	IEC104DefaultTimeout = 30
	IEC104DefaultK       = 12
	IEC104DefaultW       = 8

	// Device Status Constants
	DeviceStatusConnected    = "connected"
	DeviceStatusDisconnected = "disconnected"
	DeviceStatusError        = "error"

	// NATS Subjects
	IEC104DataSubject = "iec104.data"

	// API Routes
	IEC104BasePath       = "/iec104"
	IEC104DevicesPath    = "/devices"
	IEC104DevicePath     = "/device/:device_id"
	IEC104StatusPath     = "/status"
	IEC104DataPath       = "/data"
	IEC104ConnectPath    = "/connect"
	IEC104DisconnectPath = "/disconnect"

	// Error Messages
	ErrDeviceNotFound      = "device not found"
	ErrInvalidDeviceID     = "invalid device id"
	ErrConnectionFailed    = "connection failed"
	ErrDisconnectionFailed = "disconnection failed"
)
