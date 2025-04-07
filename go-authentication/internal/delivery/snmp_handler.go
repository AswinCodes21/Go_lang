package delivery

import (
	"net/http"
	"time"

	"go-authentication/internal/services"

	"github.com/gin-gonic/gin"
)

type SNMPHandler struct {
	snmpService *services.SNMPService
}

func NewSNMPHandler(snmpService *services.SNMPService) *SNMPHandler {
	return &SNMPHandler{
		snmpService: snmpService,
	}
}

func (h *SNMPHandler) GetDevices(c *gin.Context) {
	devices := h.snmpService.GetDevices()
	if devices == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "error",
			"error":  "SNMP simulator not initialized",
		})
		return
	}

	deviceList := make([]gin.H, 0)
	for _, device := range devices {
		// Get latest data for the device
		latestData := h.snmpService.GetLatestData(device.ID)

		deviceList = append(deviceList, gin.H{
			"id":          device.ID,
			"ip":          device.IP,
			"dataPattern": device.DataPattern,
			"oids":        device.OIDs,
			"status":      "active",
			"lastUpdate":  time.Now().Format(time.RFC3339),
			"metrics":     latestData,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"devices":   deviceList,
			"total":     len(deviceList),
			"timestamp": time.Now().Format(time.RFC3339),
		},
	})
}

func (h *SNMPHandler) GetDeviceData(c *gin.Context) {
	// TODO: Implement device data retrieval
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// HealthCheck verifies if SNMP service is running
func (h *SNMPHandler) HealthCheck(c *gin.Context) {
	running, deviceCount := h.snmpService.GetStatus()
	if !running {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "error",
			"error":  "SNMP simulator not initialized",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data": gin.H{
			"devices": deviceCount,
			"running": true,
		},
	})
}
