package main

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
)

// Main Handlers

// Sample Handler sends a default system response
func SampleHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	host := make(map[string]interface{})
	fqdn, err := os.Hostname()
	host["fqdn"] = fqdn
	if err != nil {
		log.Warn(err)
		host["fqdn"] = err
	}
	addrs, err := net.LookupHost(fqdn)
	if err != nil {
		log.Warn(err)
		host["ip_address"] = err
	}
	//addrs, err := net.InterfaceAddrs()
	//host["ip_address"] = addrs[len(addrs)-1].(*net.IPNet).IP // need something better
	host["ip_address"] = addrs[len(addrs)-1] // need something better
	data["host"] = host
	ResponseHandler(w, r, http.StatusOK, data)
}

// Hello Handler checks a parameter and returns the response
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		name = "Name does not exist"
	}
	data := make(map[string]interface{})
	sites := []string{"a", "b", "c"}
	data["method"] = r.Method
	data["url"] = fmt.Sprintf("%s", r.URL)
	data["sites"] = sites
	data["name"] = name
	ResponseHandler(w, r, http.StatusOK, data)
}
