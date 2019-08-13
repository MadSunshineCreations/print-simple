package main

import (
	"encoding/json"
	"net/http"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	var printers = dashboard.Printers

	var wg sync.WaitGroup
	wg.Add(len(Printers))
	for i := 0; i < len(printers); i++ {
		go func() {
			printers[i].getSettings()
			printers[i].getConnectionInfo()
			printers[i].getTemperatureInfo()
			wg.Done()
		}()
	}
	wg.Wait()
	// Load Job info if Printing
	for i := 0; i < len(printers); i++ {
		fillDashboardPorts(&printers[i])
		if printers[i].Connection.State == "Printing" {
			printers[i].getJobInfo()
		}
	}

	json.NewEncoder(w).Encode(dashboard)
}

func 

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
