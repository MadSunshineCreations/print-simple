package main

import (
	"encoding/json"
	"net/http"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	var printers = readConfig()
	dashboard.Printers = printers

	for i := 0; i < len(printers); i++ {
		printers[i].getConnectionInfo()
		printers[i].getTemperatureInfo()
		fillDashboardPorts(&printers[i])
	}

	json.NewEncoder(w).Encode(dashboard)
}

func fillDashboardPorts(p *Printer) {
	for i := 0; i < len(p.Connection.AvailablePorts); i++ {
		var found = -1
		for j := 0; j < len(dashboard.Ports); j++ {
			if dashboard.Ports[j].Name == p.Connection.AvailablePorts[i] {
				found = j
			}
		}
		if found == -1 {
			dashboard.Ports = append(dashboard.Ports, Port{true, p.Connection.AvailablePorts[i]})
			found = 0
		}
		if dashboard.Ports[found].Name == p.Connection.Port {
			dashboard.Ports[found].Available = false
		}
	}
}
