package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// Dashboard represents all the data for the dashboard veiew
type Dashboard struct {
	Printers []Printer `json:"printers"`
	Ports    []Port    `json:"ports"`
}

// Port represents a port and if it's in use or not
type Port struct {
	Available bool   `json:"available"`
	Name      string `json:"name"`
	HostKey   string `json:"host"`
}

// ConnectRequest represents a request to connect a printer to a port
type ConnectRequest struct {
	PrinterName string `json:"printer_name"`
	Port        string `json:"port"`
}

// PrinterRequest is for requests that only require the printer name
type PrinterRequest struct {
	PrinterName string `json:"printer_name"`
}

// JobRequest is for operations that cancel or start a print
type JobRequest struct {
	PrinterName string `json:"printer_name"`
	Operation   string `json:"operation"`
}

// MoveRequest is for moving the extruder head
type MoveRequest struct {
	PrinterName string `json:"printer_name"`
	Z           int    `json:"z"`
}

// PrintFileRequest is
type PrintFileRequest struct {
	PrinterName string `json:"printer_name"`
	FileName    string `json:"file_name"`
}

var dashboard = Dashboard{}
var dashboardMutex = &sync.Mutex{}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ready to Print!")
}

func connectHandler(w http.ResponseWriter, req *http.Request) {
	var printers = dashboard.Printers

	decoder := json.NewDecoder(req.Body)
	var r ConnectRequest
	err := decoder.Decode(&r)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(printers); i++ {
		if printers[i].Name == r.PrinterName {
			printers[i].connect(r.Port)
		}
	}

	json.NewEncoder(w).Encode(r)
}

func preheatHandler(w http.ResponseWriter, req *http.Request) {
	var printers = dashboard.Printers

	decoder := json.NewDecoder(req.Body)
	var r PrinterRequest
	err := decoder.Decode(&r)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(printers); i++ {
		if printers[i].Name == r.PrinterName {
			printers[i].preheat(200, 60)
		}
	}

	json.NewEncoder(w).Encode(r)
}

func extrudeHandler(w http.ResponseWriter, req *http.Request) {
	var printers = dashboard.Printers

	decoder := json.NewDecoder(req.Body)
	var r PrinterRequest
	err := decoder.Decode(&r)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(printers); i++ {
		if printers[i].Name == r.PrinterName {
			printers[i].extrude(100)
		}
	}

	json.NewEncoder(w).Encode(r)
}

func jobHandler(w http.ResponseWriter, req *http.Request) {
	var printers = dashboard.Printers

	decoder := json.NewDecoder(req.Body)
	var r JobRequest
	err := decoder.Decode(&r)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(printers); i++ {
		if printers[i].Name == r.PrinterName {
			if r.Operation == "cancel" {
				printers[i].cancel()
				printers[i].preheat(0, 0)
			} else if r.Operation == "start" {
				printers[i].start()
			}
		}
	}

	json.NewEncoder(w).Encode(r)
}

func movezHandler(w http.ResponseWriter, req *http.Request) {
	var printers = dashboard.Printers

	decoder := json.NewDecoder(req.Body)
	var r MoveRequest
	err := decoder.Decode(&r)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(printers); i++ {
		if printers[i].Name == r.PrinterName {
			printers[i].movez(r.Z)
		}
	}

	json.NewEncoder(w).Encode(r)
}

func printFileHandler(w http.ResponseWriter, req *http.Request) {
	var printers = dashboard.Printers

	decoder := json.NewDecoder(req.Body)
	var r PrintFileRequest
	err := decoder.Decode(&r)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(printers); i++ {
		if printers[i].Name == r.PrinterName {
			printers[i].printFile(r.FileName)
		}
	}

	json.NewEncoder(w).Encode(r)
}

func readConfig() []Printer {
	data, _ := ioutil.ReadFile("printer-list.yaml")

	printers := []Printer{}
	yaml.Unmarshal([]byte(data), &printers)

	return printers
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Printf("Loading Printers from Config\n")
	var printers = readConfig()
	dashboard.Printers = printers
	//Start status loop
	go func() {
		for {
			loadStatus()
			time.Sleep(time.Second * 1)
		}
	}()

	fmt.Printf("Loaded!\n")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/status/", statusHandler)
	http.HandleFunc("/connect", connectHandler)
	http.HandleFunc("/connect/", connectHandler)
	http.HandleFunc("/preheat", preheatHandler)
	http.HandleFunc("/preheat/", preheatHandler)
	http.HandleFunc("/extrude", extrudeHandler)
	http.HandleFunc("/extrude/", extrudeHandler)
	http.HandleFunc("/job", jobHandler)
	http.HandleFunc("/job/", jobHandler)
	http.HandleFunc("/movez", movezHandler)
	http.HandleFunc("/movez/", movezHandler)
	http.HandleFunc("/print_file/", printFileHandler)
	http.HandleFunc("/print_file", printFileHandler)
	fmt.Printf("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
