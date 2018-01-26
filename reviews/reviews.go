/*
Reviews package contains the abstraction layer into the Beer Reviews service
*/
package main

import (
	"time"
)

/*
Review is a structure used for serializing/deserializing data for the database.
It is comprised of additional object dependencies: Beer and Reviewer.

    {
        "text" : "I had this with a burger. Excellent",
        "beer": { "id": 1 },
        "reviewer": { "id": 1, "name": "chris" }
    }

*/
type Review struct {
	Id        string `json:"id"`
	Text      string `json:"text,omitempty"`
	Beer      `json:"beer"`
	Reviewer  `json:"reviewer"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	//Tags     []string              `json:"tags,omitempty"`
}

// Beer struct
type Beer struct {
	Id int `json:"id"`
}

// Reviewer struct
type Reviewer struct {
	Id   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

/*
Index is a structure used for serializing/deserializing data in Elasticsearch.
    {
        "index" :
            {   "_index" : "beer",
                "_type" : "review",
                "_id" : "da8d2d25-4069-478b-8737-fe2312c42352"
             }
    }

*/
type Index struct {
	Index struct {
		Name string `json:"_index"`
		Type string `json:"_type"`
		Id   string `json:"_id"`
	}
}

/*
Review index is both Index and Review structu
*/
type ReviewIndex struct {
	Index
	Review
}
