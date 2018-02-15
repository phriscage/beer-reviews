package main

import (
	"fmt"
	//"github.com/gorilla/mux"
	"github.com/fatih/structs"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"net"
	//"encoding/json"
	"net/http"
	"os"
)

type DebugInfo struct {
	Host    interface{} `structs:"host,omitempty"`
	Request interface{} `structs:"request,omitempty"`
}

// Get some host and HTTP request information and structure the data
func GetDebugInfo(r *http.Request) *DebugInfo {
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

	request := make(map[string]interface{})
	request["url"] = fmt.Sprintf("%s%s", r.Host, r.URL.Path)
	request["headers"] = r.Header

	d := DebugInfo{
		Host:    host,
		Request: request,
	}
	return &d
}

// Health Handler returns a heathcheck for the service
func (env *Env) HealthHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	debug := vals["debug"]
	//log.Debug(vals)

	// Ping the Elasticsearch server to get e.g. the version number
	_, _, err := env.client.Ping(env.elasticUrl).Do(ctx)
	if err != nil {
		// Handle error
		log.Warn(err)
		if e, ok := err.(*elastic.Error); ok {
			if e.Status == 404 {
				NotFoundHandler(w, r)
				return
			}
		}
		ResponseErrorHandler(w, r, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	if debug != nil {
		data := structs.Map(GetDebugInfo(r))
		//log.Debug(data)
		ResponseHandler(w, r, http.StatusOK, data)
	} else {
		ResponseHandler(w, r, http.StatusOK, nil)
	}
}
