package main

import (
	"fmt"

	octoprint "github.com/mcuadros/go-octoprint"
)

//Printer represents a octoprint Printer name and config
type Printer struct {
	Name       string `yaml:"name" json:"name"`
	URL        string `yaml:"url" json:"-"`
	Apikey     string `yaml:"api_key" json:"-"`
	Connection struct {
		State          string   `json:"state"`
		Port           string   `json:"port"`
		AvailablePorts []string `json:"-"`
	} `json:"connection"`
	Temperature struct {
		Hotend octoprint.TemperatureData `json:"hotend"`
		Bed    octoprint.TemperatureData `json:"bed"`
	} `json:"temperature"`
}

func (p *Printer) getConnectionInfo() {
	client := octoprint.NewClient(p.URL, p.Apikey)
	req := octoprint.ConnectionRequest{}
	s, _ := req.Do(client)
	p.Connection.State = string(s.Current.State)
	p.Connection.Port = s.Current.Port
	p.Connection.AvailablePorts = s.Options.Ports
}

func (p *Printer) getTemperatureInfo() {
	client := octoprint.NewClient(p.URL, p.Apikey)
	req := octoprint.StateRequest{}
	s, err := req.Do(client)
	fmt.Printf("Temp Response %v", p)
	if err == nil {
		p.Temperature.Hotend = s.Temperature.Current["tool0"]
		p.Temperature.Bed = s.Temperature.Current["bed"]
	}
}

func (p *Printer) connect() {
	client := octoprint.NewClient("https://ender32.madsunshinecreations.com/", "9EDCA3B52DFC4B2AB18D6E8616E2D31B")
	r := octoprint.ConnectionRequest{}

	s, _ := r.Do(client)

	fmt.Printf("Connection State: %q\n", s.Current.State)
	if s.Current.State == "Closed" {
		req := octoprint.ConnectRequest{}
		req.Do(client)
	}
}
