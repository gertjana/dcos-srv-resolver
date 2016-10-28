package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type Response struct {
	Name    string   `json:"name"`
	Service string   `json:"service"`
	Targets []Target `json:"targets"`
}

type Target struct {
	Host     string   `json:"host"`
	Ips      []net.IP `json:"ips"`
	Port     uint16   `json:"port"`
	Priority uint16   `json:"priority"`
	Weight   uint16   `json:"weight"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func createResponse(name string, service string, srvs []*net.SRV) string {
	targets := make([]Target, len(srvs), (cap(srvs)+1)*2)
	for i := range srvs {
		ips, err := net.LookupIP(srvs[i].Target)
		if err != nil {
			log.Println(err.Error())
		}
		targets[i] = Target{
			Host:     srvs[i].Target,
			Ips:      ips,
			Port:     srvs[i].Port,
			Priority: srvs[i].Priority,
			Weight:   srvs[i].Weight,
		}
	}
	response := &Response{
		Name:    name,
		Service: service,
		Targets: targets,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println(err.Error())
	}
	return string(jsonResponse)
}

func lookup(w http.ResponseWriter, r *http.Request) {
	host := "marathon.mesos."
	const protocol = "tcp"

	if strings.HasPrefix(r.Host, "localhost") {
		log.Println("localhost detected, using test SRV records at addictivesoftware.net")
		host = "addictivesoftware.net."
	}

	service := mux.Vars(r)["service"]

	w.Header().Set("Content-Type", "application/json")

	if cname, srvs, err := net.LookupSRV(service, protocol, host); err != nil {
		errorResponse, _ := json.Marshal(&ErrorResponse{
			Error: err.Error(),
		})
		w.WriteHeader(404)
		io.WriteString(w, string(errorResponse))
	} else {
		io.WriteString(w, createResponse(cname, service, srvs))
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
}

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/service/{service}", lookup).Methods("GET")
	rtr.HandleFunc("/status", status).Methods("GET")

	http.Handle("/", rtr)
	http.ListenAndServe(":8000", nil)
}
