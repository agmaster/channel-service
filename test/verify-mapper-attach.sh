## 
##  Test commands for mapper attachment plugin in Elasticsearch
##
##


# Create a property mapping using the new type attachment:
curl -XPOST  -H 'Content-Type: application/json' -d \
'{
  "mappings": {
    "person": {
      "properties": {
        "cv": { "type": "attachment" }
}}}}'  -v http://localhost:9200/try-attachments


#Index a new document populated with a base64-encoded attachment:

curl -XPOST  -H 'Content-Type: application/json' -d \
'{
  "cv": "e1xydGYxXGFuc2kNCkxvcmVtIGlwc3VtIGRvbG9yIHNpdCBhbWV0DQpccGFyIH0="
}'   -v  http://localhost:9200/try-attachments/person/1


#Search for the document using words in the attachment:
# If you get a hit for your indexed document, the plugin should be installed and working.

curl -XPOST  -H 'Content-Type: application/json' -d \
'{
  "query": {
    "query_string": {
      "query": "ipsum"
}}}'   -v  http://localhost:9200/try-attachments/person/_search


# Usage
curl -XDELETE  http://localhost:9200/test
curl -XPUT  http://localhost:9200/test

curl -XPUT -v  http://localhost:9200/test/person/_mapping -d \
 '{
    "person" : {
        "properties" : {
            "my_attachment" : { "type" : "attachment" }
        }
    }
}'

# Json index 
# if content type, resource name or language need
curl -XPUT   -d \
'{
    "attachment" :  "... base64 encoded attachment ..."
}'  -v   http://localhost:9200/test/person/3

curl -XPOST  -H 'Content-Type: application/json' -d '{
  "cv": "e1xydGYxXGFuc2kNCkxvcmVtIGlwc3VtIGRvbG9yIHNpdCBhbWV0DQpccGFyIH0="
}'   -v http://localhost:9200/test/person/1 


curl -XPOST  -H 'Content-Type: application/json' -d \
'{

    "mypdf" : {
        "_content_type" : "application/pdf",
        "_name" : "resource/name/of/my.pdf",
        "_language" : "en",
        "_content" : "... base64 encoded attachment ..."
    }
}'  -v http://localhost:9200/test/person/3


curl -XGET -v http://localhost:9200/test/person/3


curl -XPUT  -H 'Content-Type: application/json'  -v  http://localhost:9200/test/person/4 -d \
'{
    "myattachment4" : {
        "_content_type" : "application/pdf",
        "_name" : "resource/name/of/my.pdf",
        "_language" : "en",
        "_content" : "e1xydGYxXGFuc2kNCkxvcmVtIGlwc3VtIGRvbG9yIHNpdCBhbWV0DQpccGFyIH0"
    }
}'  

#Querying or accessing metadata
curl -XDELETE  http://localhost:9200/metadata
curl -XPUT  http://localhost:9200/metadata

curl -XPUT http://localhost:9200/metadata/person/_mapping -d \
'{
  "person": {
    "properties": {
      "file": {
        "type": "attachment",
        "fields": {
          "content_type": {
            "type": "string",
            "store": true
          }
        }
      }
    }
  }
}'


#PUT /test/person/1?refresh=true
curl -XPUT http://localhost:9200/metadata/person/1?refresh=true -d \
'{
  "file": "IkdvZCBTYXZlIHRoZSBRdWVlbiIgKGFsdGVybmF0aXZlbHkgIkdvZCBTYXZlIHRoZSBLaW5nIg=="
}'


#GET /test/person/_search
curl -XGET http://localhost:9200/metadata/person/_search -d \
'{
  "fields": [ "file.content_type" ],
  "query": {
    "match": {
      "file.content_type": "text plain"
    }
  }
}'
