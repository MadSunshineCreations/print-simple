package main

import (
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

var dashboard = Dashboard{}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ready to Print!")
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
	fmt.Printf("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
