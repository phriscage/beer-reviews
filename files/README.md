## Create Index

        curl -XPUT 'localhost:9200/beer?pretty' -H 'Content-Type: application/json' -d @index.json

## Create Entries

        curl -s -H "Content-Type: application/x-ndjson" -XPOST localhost:9200/_bulk --data-binary "@seed.json"; echo

## Query entries

        curl -s http://localhost:9200/beer/_search?q=beer.id:1
