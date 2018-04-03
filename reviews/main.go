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

const (
	elasticBeerIndex = "beers"
	elasticBeerType  = "review"
)

// Global environment struct for handling the database connection
type Env struct {
	client     *elastic.Client
	elasticUrl string
}

var (
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx = context.Background()
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

// init database for environment struct
func initDb(elasticUrl string) (*elastic.Client, error) {
	// Obtain a client and connect to the default Elasticsearch installation
	c, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(elasticUrl),
		elastic.SetHealthcheckTimeoutStartup(10*time.Second),
	)
	if err != nil {
		// Handle error
		return nil, err
	}
	defer c.Stop()

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := c.Ping(elasticUrl).Do(ctx)
	if err != nil {
		// Handle error
		return nil, err
	}
	log.Debugf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := c.ElasticsearchVersion(elasticUrl)
	if err != nil {
		// Handle error
		return nil, err
	}
	log.Infof("Elasticsearch version %s\n", esversion)
	return c, nil
}

// main logic
func main() {

	elasticUrl := os.Getenv("ELASTICSEARCH_DATABASE_URI")
	log.Debugf("ELASTICSEARCH_DATABASE_URI: %s", os.Getenv("ELASTICSEARCH_DATABASE_URI"))
	if elasticUrl == "" {
		elasticUrl = "http://localhost:9200"
	}
	c, err := initDb(elasticUrl)
	if err != nil {
		log.Fatal(err)
	}

	exists, err := c.IndexExists(elasticBeerIndex).Do(ctx)
	if err != nil {
		// Handle error
		log.Fatal(err)
		panic(err)
	}
	if !exists {
		// Index does not exist yet.
		log.Infof("'%s' index does not exist yet, seeding data...", elasticBeerIndex)
		err := seedData(c)
		if err != nil {
			log.Fatal(err)
		}
	}
	env := &Env{client: c, elasticUrl: elasticUrl}

	// Main router logic
	addr := flag.String("addr", ":8080", "http listen address")
	flag.Parse()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/sample", SampleHandler).Methods("GET")
	router.HandleFunc("/sample", MethodNotAllowedHandler)
	router.HandleFunc("/reviews", env.ReviewsHandler).Methods("GET")
	router.HandleFunc("/reviews", MethodNotAllowedHandler)
	router.HandleFunc("/review/{id}", env.ReviewsIdHandler).Methods("GET")
	router.HandleFunc("/review/{id}", MethodNotAllowedHandler)
	router.HandleFunc("/reviewer/{id}/reviews", env.ReviewersIdReviewsHandler).Methods("GET")
	router.HandleFunc("/reviewer/{id}/reviews", MethodNotAllowedHandler)
	router.HandleFunc("/beer/{id}/reviews", env.BeersIdReviewsHandler).Methods("GET")
	router.HandleFunc("/beer/{id}/reviews", MethodNotAllowedHandler)
	router.HandleFunc("/health", env.HealthHandler).Methods("GET")
	router.HandleFunc("/health", MethodNotAllowedHandler)
	//router.HandleFunc("/blinkts/{action}/{id}", BlinktsHandler).Methods("POST")
	//router.Handle("/blinkts/random", handlers.MethodHandler{
	//"POST": http.HandlerFunc(BlinktsHandler),
	//})
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	log.Infof("Starting web server on " + *addr)
	log.Fatal(http.ListenAndServe(*addr, handlers.CombinedLoggingHandler(os.Stderr, router)))
}
