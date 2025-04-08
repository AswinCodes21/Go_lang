package delivery

import (
	"net/http"

	"go-authentication/internal/domain"
	"go-authentication/internal/services"

	"github.com/gin-gonic/gin"
)

type IEC104Handler struct {
	service *services.IEC104Service
}

func NewIEC104Handler(service *services.IEC104Service) *IEC104Handler {
	return &IEC104Handler{
		service: service,
	}
}

// GetDevices returns a list of all IEC104 devices
func (h *IEC104Handler) GetDevices(c *gin.Context) {
	devices := h.service.GetDevices()
	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
	})
}

// GetDeviceStatus returns the status of a specific device
func (h *IEC104Handler) GetDeviceStatus(c *gin.Context) {
	deviceID := c.Param("device_id")
	status := h.service.GetDeviceStatus(deviceID)
	c.JSON(http.StatusOK, gin.H{
		"device_id": deviceID,
		"status":    status,
	})
}

// GetLatestData returns the latest data for a specific device
func (h *IEC104Handler) GetLatestData(c *gin.Context) {
	deviceID := c.Param("device_id")
	data := h.service.GetLatestData(deviceID)
	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": domain.ErrDeviceNotFound,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"device_id": deviceID,
		"data":      data,
	})
}

// ConnectToDevice initiates a connection to a specific device
func (h *IEC104Handler) ConnectToDevice(c *gin.Context) {
	deviceID := c.Param("device_id")
	h.service.ConnectToDevice(deviceID)
	c.JSON(http.StatusOK, gin.H{
		"message":   "Connection initiated",
		"device_id": deviceID,
	})
}

// DisconnectFromDevice disconnects from a specific device
func (h *IEC104Handler) DisconnectFromDevice(c *gin.Context) {
	deviceID := c.Param("device_id")
	h.service.DisconnectFromDevice(deviceID)
	c.JSON(http.StatusOK, gin.H{
		"message":   "Disconnection initiated",
		"device_id": deviceID,
	})
}
