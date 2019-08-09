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
}

var dashboard = Dashboard{}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func readConfig() []Printer {
	data, _ := ioutil.ReadFile("printer-list.yaml")

	printers := []Printer{}
	yaml.Unmarshal([]byte(data), &printers)

	fmt.Printf("--- Printers:\n%v\n\n", printers)
	return printers
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/status/", statusHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
