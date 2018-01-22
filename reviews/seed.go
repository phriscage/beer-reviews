package main

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"os"
)

var seedDataFile = "data/seed.json"

// seed the sample data from
func seedData(client *elastic.Client) error {
	// Open our jsonFile
	jsonFile, err := os.Open(seedDataFile)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	jsonScanner := bufio.NewScanner(jsonFile)

	count := 0
	index := Index{}
	for jsonScanner.Scan() {
		if count%2 == 0 {
			//index = Index{}
			err := json.Unmarshal([]byte(jsonScanner.Text()), &index)
			if err != nil {
				return err
			}
		} else {
			review := Review{}
			err := json.Unmarshal([]byte(jsonScanner.Text()), &review)
			if err != nil {
				return err
			}
			log.Debugf("Beer: %+v, Reviewer: %+v, Index: %+v\n", review.Beer, review.Reviewer, index.Index)
			indexReq := elastic.NewBulkIndexRequest().Index(index.Index.Name).Type(index.Index.Type).Id(index.Index.Id).Doc(review)
			bulkRequest := client.Bulk()
			bulkRequest = bulkRequest.Add(indexReq)

			// Do sends the bulk requests to Elasticsearch
			bulkResponse, err := bulkRequest.Do(context.Background())
			if err != nil {
				// ...
				panic(err)
			}
			// Failed() returns information about failed bulk requests
			// (those with a HTTP status code outside [200,299].
			failedResults := bulkResponse.Failed()
			if failedResults != nil {
				// ...
				panic(failedResults)
			}
		}
		count++
	}
	return nil
}
