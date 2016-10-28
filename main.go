package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net"
	"net/http"
	"strings"
)

type Response struct {
	Name    string
	Service string
	Targets []Target
}

type Target struct {
	Host     string
	Ips      []net.IP
	Port     uint16
	Priority uint16
	Weight   uint16
}

type ErrorResponse struct {
	Error string
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func CreateResponse(name string, service string, srvs []*net.SRV) string {
	targets := make([]Target, len(srvs), (cap(srvs)+1)*2)
	for i := range srvs {
		ips, _ := net.LookupIP(srvs[i].Target)
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
	jsonResponse, _ := json.Marshal(response)
	return string(jsonResponse)
}

func lookup(w http.ResponseWriter, r *http.Request) {
	host := "marathon.mesos."
	const protocol = "tcp"

	// to locally test this i've added 2 srv records for _test._tcp to the addictive software domain
	if strings.HasPrefix(r.Host, "localhost") {
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
		io.WriteString(w, CreateResponse(cname, service, srvs))
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
