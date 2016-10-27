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
	Service string     `json:"Service"`
	Srvs    []*net.SRV `json:"Targets"`
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

func ToJson(service string, srvs []*net.SRV) string {
	response := &Response{
		Service: service,
		Srvs:    srvs,
	}
	jsonResponse, _ := json.Marshal(response)
	return string(jsonResponse)
}

func lookup(w http.ResponseWriter, r *http.Request) {
	var host = "marathon.mesos."
	if strings.HasPrefix(r.Host, "localhost") { // to locally test this i've added 2 srv records to my domain
		host = "addictivesoftware.net."
	}

	service := mux.Vars(r)["service"]

	w.Header().Set("Content-Type", "application/json")
	if _, srvs, err := net.LookupSRV(service, "tcp", host); err != nil {
		errorResponse, _ := json.Marshal(&ErrorResponse{
			Error: err.Error(),
		})
		io.WriteString(w, string(errorResponse))
	} else {
		io.WriteString(w, ToJson(service, srvs))
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
