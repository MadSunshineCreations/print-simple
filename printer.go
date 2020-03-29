package main

import (
	"fmt"

	octoprint "github.com/mark579/go-octoprint"
)

//Printer represents a octoprint Printer name and config
type Printer struct {
	Name       string `yaml:"name" json:"name"`
	URL        string `yaml:"url" json:"-"`
	Apikey     string `yaml:"api_key" json:"-"`
	HostKey    string `yaml:"host_key" json:"-"`
	Connection struct {
		State          string   `json:"state"`
		Port           string   `json:"port"`
		AvailablePorts []string `json:"-"`
	} `json:"connection"`
	Temperature struct {
		Hotend octoprint.TemperatureData `json:"hotend"`
		Bed    octoprint.TemperatureData `json:"bed"`
	} `json:"temperature"`
	Job struct {
		Name     string  `json:"name"`
		Progress float64 `json:"progress"`
	} `json:"job"`
	Settings struct {
		Loaded          bool   `json:"-"`
		WebcamStreamURL string `json:"webcam_stream_url"`
	} `json:"settings"`
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

	if err == nil {
		p.Temperature.Hotend = s.Temperature.Current["tool0"]
		p.Temperature.Bed = s.Temperature.Current["bed"]
	}
}

func (p *Printer) getJobInfo() {
	client := octoprint.NewClient(p.URL, p.Apikey)
	req := octoprint.JobRequest{}
	rep, err := req.Do(client)

	if err == nil {
		p.Job.Name = rep.Job.File.Name
		p.Job.Progress = rep.Progress.Completion
	}
}

func (p *Printer) getSettings() {
	if p.Settings.Loaded == true {
		return
	}
	client := octoprint.NewClient(p.URL, p.Apikey)
	req := octoprint.SettingsRequest{}
	rep, err := req.Do(client)

	if err == nil {
		p.Settings.WebcamStreamURL = rep.Webcam.StreamURL
		p.Settings.Loaded = true
	} else {
		fmt.Printf("%+v\n", err)
	}
}

func (p *Printer) connect(port string) {
	client := octoprint.NewClient(p.URL, p.Apikey)
	r := octoprint.ConnectionRequest{}
	s, _ := r.Do(client)

	if s.Current.State == "Closed" {
		req := octoprint.ConnectRequest{}
		req.Port = port
		req.Do(client)
	}
}

func (p *Printer) preheat(tool float64, bed float64) {
	var bedTarget float64
	client := octoprint.NewClient(p.URL, p.Apikey)

	target := make(map[string]float64)

	target["tool0"] = tool
	bedTarget = bed
	toolRequest := octoprint.ToolTargetRequest{target}
	bedRequest := octoprint.BedTargetRequest{bedTarget}
	toolRequest.Do(client)
	bedRequest.Do(client)
}

func (p *Printer) extrude(amount int) {
	client := octoprint.NewClient(p.URL, p.Apikey)

	r := octoprint.ToolExtrudeRequest{amount}
	r.Do(client)
}

func (p *Printer) cancel() {
	client := octoprint.NewClient(p.URL, p.Apikey)

	r := octoprint.CancelRequest{}
	r.Do(client)
}

func (p *Printer) start() {
	client := octoprint.NewClient(p.URL, p.Apikey)

	r := octoprint.StartRequest{}
	r.Do(client)
}

func (p *Printer) movez(amount int) {
	client := octoprint.NewClient(p.URL, p.Apikey)

	r := octoprint.PrintHeadJogRequest{0, 0, amount, false, 200}
	r.Do(client)
}
