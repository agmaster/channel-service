#!/bin/sh

curl -XDELETE 'http://localhost:9200/test'
curl -XPUT 'http://localhost:9200/test'
curl -XPUT 'http://localhost:9200/test/person/_mapping' -d \
'{
  "person": {
    "properties": {
      "file": {
        "type": "attachment",
        "fields": {
          "content": {
            "type": "string",
            "term_vector":"with_positions_offsets",
            "store": true
          }
        }
      }
    }
  }
}'

curl -XPUT 'http://localhost:9200/test/person/1?refresh=true' -d \
'{
  "file": "IkdvZCBTYXZlIHRoZSBRdWVlbiIgKGFsdGVybmF0aXZlbHkgIkdvZCBTYXZlIHRoZSBLaW5nIg=="
}'


curl -XGET 'http://localhost:9200/test/person/_search?pretty=true' -d \
'{
  "fields": [],
  "query": {
    "match": {
      "file.content": "king queen"
    }
  },
  "highlight": {
    "fields": {
      "file.content": {
      }
    }
  }
}'
