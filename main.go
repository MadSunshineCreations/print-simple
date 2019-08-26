package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

var dashboard = Dashboard{}

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
	log.Println(r.PrinterName)

	for i := 0; i < len(printers); i++ {
		if printers[i].Name == r.PrinterName {
			printers[i].connect(r.Port)
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

func main() {
	fmt.Printf("Loading Printers from Config\n")
	var printers = readConfig()
	dashboard.Printers = printers
	fmt.Printf("Loaded!\n")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/status/", statusHandler)
	http.HandleFunc("/connect/", connectHandler)
	fmt.Printf("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
