package main

import (
	"context"
	"flag"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

// Init
func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)
}

// main logic
func main() {

	log.Debugf("ELASTICSEARCH_DATABASE_URI: %s", os.Getenv("ELASTICSEARCH_DATABASE_URI"))
	elasticUrl := os.Getenv("ELASTICSEARCH_DATABASE_URI")
	if elasticUrl == "" {
		elasticUrl = "http://localhost:9200"
	}
	elasticBeerIndex := "beers"

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(elasticUrl),
		elastic.SetHealthcheckTimeoutStartup(10*time.Second),
	)
	if err != nil {
		// Handle error
		log.Fatal(err)
		panic(err)
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping(elasticUrl).Do(ctx)
	if err != nil {
		// Handle error
		log.Fatal(err)
		panic(err)
	}
	log.Debugf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion(elasticUrl)
	if err != nil {
		// Handle error
		log.Fatal(err)
		panic(err)
	}
	log.Infof("Elasticsearch version %s\n", esversion)

	exists, err := client.IndexExists(elasticBeerIndex).Do(context.Background())
	if err != nil {
		// Handle error
		log.Fatal(err)
		panic(err)
	}
	if !exists {
		// Index does not exist yet.
		log.Infof("'%s' index does not exist yet, seeding data...", elasticBeerIndex)
		err := seedData(client)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Main router logic
	addr := flag.String("addr", ":8080", "http listen address")
	flag.Parse()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/sample", SampleHandler).Methods("GET")
	router.HandleFunc("/sample", MethodNotAllowedHandler)
	router.HandleFunc("/hello", HelloHandler).Methods("GET")
	router.HandleFunc("/hello", MethodNotAllowedHandler)
	router.HandleFunc("/hello/{name}", HelloHandler).Methods("GET")
	router.HandleFunc("/hello/{name}", MethodNotAllowedHandler)
	//router.HandleFunc("/blinkts/{action}", BlinktsHandler).Methods("POST")
	//router.HandleFunc("/blinkts/{action}", MethodNotAllowedHandler)
	//router.HandleFunc("/blinkts/{action}/{id}", BlinktsHandler).Methods("POST")
	//router.HandleFunc("/blinkts/{action}/{id}", MethodNotAllowedHandler)
	//router.Handle("/blinkts/random", handlers.MethodHandler{
	//"POST": http.HandlerFunc(BlinktsHandler),
	//})
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	log.Infof("Starting web server on " + *addr)
	log.Fatal(http.ListenAndServe(*addr, handlers.CombinedLoggingHandler(os.Stderr, router)))
}
