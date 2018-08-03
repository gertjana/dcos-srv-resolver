package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
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

func createShortResponse(name string, service string, srvs []*net.SRV) string {
	targets := ""
	for i := range srvs {
		target := srvs[i].Target
		port := srvs[i].Port
		ips, _ := net.LookupIP(srvs[i].Target)
		if len(ips) == 0 {
			targets += target + ":" + strconv.Itoa(int(port)) + "\n"
		} else {
			targets += ips[0].String() + ":" + strconv.Itoa(int(port)) + "\n"
		}
	}
	return targets
}

func createLongResponse(name string, service string, srvs []*net.SRV) string {
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

func lookupShort(w http.ResponseWriter, r *http.Request) {
	lookup(w, r, true)
}

func lookupLong(w http.ResponseWriter, r *http.Request) {
	lookup(w, r, false)
}

func lookup(w http.ResponseWriter, r *http.Request, short bool) {
	host := "marathon.mesos."
	const protocol = "tcp"

	if strings.HasPrefix(r.Host, "localhost") {
		log.Println("localhost detected, using test SRV records at addictivesoftware.net")
		host = "addictivesoftware.net."
	}

	service := mux.Vars(r)["service"]

	if cname, srvs, err := net.LookupSRV(service, protocol, host); err != nil {
		errorResponse, _ := json.Marshal(&ErrorResponse{
			Error: err.Error(),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		io.WriteString(w, string(errorResponse))
	} else {
		if short {
			w.Header().Set("Content-Type", "text/plain")
			io.WriteString(w, createShortResponse(cname, service, srvs))
		} else {
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, createLongResponse(cname, service, srvs))
		}
	}
}

func headers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("Header field %q, Value %q\n", k, v))
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
}

func info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, `
		<!DOCTYPE html>
		<html lang="en">
		  <head></head>
		  <body>
		  	<h1>Yellow Pages</h1>
		  	<p>Find out information about the environment this app is running in</p>
		  	<ul>
		  		<li><a href="/">This Page</a></li>
		  		<li><a href="/headers">Http Headers</a></li>
		  		<li><a href="/service/yellowpages">DNS SRV records for the names application</a></li>
		  		<li><a href="/short/yellowpages">DNS SRV Records ip:port only</a></li>
		  		<li><a href="/status">Healthcheck</a></li>
		  	</ul>
		  	<p> Source: <a href="https://github.com/gertjana/dcos-srv-resolver">https://github.com/gertjana/dcos-srv-resolver</a>
		  </body>

		</html>`)
}

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/service/{service}", lookupLong).Methods("GET")
	rtr.HandleFunc("/short/{service}", lookupShort).Methods("GET")
	rtr.HandleFunc("/status", status).Methods("GET")
	rtr.HandleFunc("/headers", headers).Methods("GET")
	rtr.HandleFunc("/", info).Methods("GET")

	http.Handle("/", rtr)
	http.ListenAndServe(":8000", nil)
}
