package main

import (
	"fmt"
	"os"
	"path/filepath"

	octoprint "github.com/mark579/go-octoprint"
)

//Printer represents a octoprint Printer name and config
type Printer struct {
	Name       string `yaml:"name" json:"name"`
	URL        string `yaml:"url" json:"-"`
	Apikey     string `yaml:"api_key" json:"-"`
	HostKey    string `yaml:"host_key" json:"-"`
	GcodeDir   string `yaml:"gcode_dir" json:"-"`
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
	Files []string `json:"files"`
}

func (p *Printer) getConnectionInfo() {
	client := octoprint.NewClient(p.URL, p.Apikey)
	req := octoprint.ConnectionRequest{}
	s, _ := req.Do(client)
	p.Connection.State = string(s.Current.State)
	p.Connection.Port = s.Current.Port
	p.Connection.AvailablePorts = s.Options.Ports
}

func (p *Printer) loadGcodeFiles() {
	p.Files = nil
	err := filepath.Walk(p.GcodeDir, func(path string, info os.FileInfo, err error) error {
		if path != p.GcodeDir {
			p.Files = append(p.Files, path)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Unable to load GCode Files")
	}
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

func (p *Printer) printFile(fileName string) {
	client := octoprint.NewClient(p.URL, p.Apikey)

	r := octoprint.UploadFileRequest{Location: octoprint.Local, Select: true, Print: true}
	file, _ := os.Open(fileName)
	r.AddFile(file.Name(), file)
	r.Do(client)
}
