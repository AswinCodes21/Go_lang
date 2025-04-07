package snmp

import (
	"os"
	"time"
)

func LoadConfig() *SNMPConfig {
	community := os.Getenv("SNMP_COMMUNITY")
	if community == "" {
		community = "public"
	}

	interval := os.Getenv("SNMP_INTERVAL")
	if interval == "" {
		interval = "30s"
	}

	simInterval, _ := time.ParseDuration(interval)

	return &SNMPConfig{
		Community: community,
		Port:      161,
		Interval:  simInterval,
		Devices: []MockDevice{
			{
				ID:          "device1",
				IP:          "192.168.1.1",
				Community:   community,
				OIDs:        []string{"1.3.6.1.2.1.1.1.0", "1.3.6.1.2.1.1.3.0"},
				DataPattern: "random",
			},
		},
	}
}
