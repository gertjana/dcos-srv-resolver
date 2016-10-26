package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net"
	"net/http"
	"strings"
)

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func lookup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	if _, srvs, err := net.LookupSRV(service, "tcp", "marathon.mesos"); err != nil {
		io.WriteString(w, fmt.Sprintf("ERROR:%s", err))
	} else {
		io.WriteString(w, fmt.Sprintf("%s:%d", TrimSuffix(srvs[0].Target, "."), srvs[0].Port))
	}
}

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/{service}", lookup).Methods("GET")

	http.Handle("/", rtr)
	http.ListenAndServe(":8000", nil)
}
