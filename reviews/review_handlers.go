package main

import (
	//"fmt"
	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	//"net"
	"encoding/json"
	"net/http"
	//"os"
)

// Reviews Handler returns all the reviews for a given reviewer_id
func ReviewsHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	searchResult, err := client.Search().
		Index(elasticBeerIndex).
		//Query(termQuery).   // specify the query
		//Sort("user", true). // sort by "user" field, ascending
		From(0).Size(10). // take documents 0-9
		Pretty(true).     // pretty print request and response JSON
		Do(ctx)           // execute
	if err != nil {
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
	log.Debug(searchResult.Hits.Hits)
	// Iterate through results
	reviews := []Review{}
	for _, hit := range searchResult.Hits.Hits {
		// hit.Index contains the name of the index
		// Deserialize hit.Source into a Review (could also be just a map[string]interface{}).
		review := Review{}
		err := json.Unmarshal(*hit.Source, &review)
		if err != nil {
			// Deserialization failed
			log.Warn(err)
			ResponseErrorHandler(w, r, http.StatusInternalServerError, []string{err.Error()})
			return
		}
		reviews = append(reviews, review)
	}
	data["reviews"] = reviews
	ResponseHandler(w, r, http.StatusOK, data)
}

// Reviews Handler returns all the reviews for a given reviewer_id
func ReviewsIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug(vars)
	data := make(map[string]interface{})
	if _, ok := vars["id"]; !ok {
		NotFoundHandler(w, r)
		return
	}
	review, err := client.Get().
		Index(elasticBeerIndex).
		Type(elasticBeerType).
		Id(vars["id"]).
		Do(ctx)
	if err != nil {
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
	log.Debug(review)
	if review != nil && review.Found {
		data["review"] = review.Source
		ResponseHandler(w, r, http.StatusOK, data)
	} else {
		NotFoundHandler(w, r)
		return
	}
}
