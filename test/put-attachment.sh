#!/bin/sh
# Example of handling attachments in Elasticsearch
# ================================================
#
# 1. Install the plugin: $ plugin -install elasticsearch/elasticsearch-mapper-attachments/1.7.0
# 2. Run the script:     $ bash test_attachments.sh
#
# More info: <http://www.elasticsearch.org/guide/reference/mapping/attachment-type/>

curl -X DELETE http://localhost:9200/test_attachments

echo;
curl -X POST http://localhost:9200/test_attachments -d '{
  "mappings" : {
    "document" : {
      "properties" : {
        "content" : {
          "type" : "attachment",
          "fields" : {
            "content"  : { "store" : "yes" },
            "author"   : { "store" : "yes" },
            "title"    : { "store" : "yes", "analyzer" : "english"},
            "date"     : { "store" : "yes" },
            "keywords" : { "store" : "yes", "analyzer" : "keyword" },
            "_name"    : { "store" : "yes" },
            "content_length" : { "store" : "yes" },
            "content_type" : { "store" : "yes" }
          }
        }
      }
    }
  }
}'

echo;
echo '>>> Index the document'
curl -i -X PUT http://localhost:9200/test_attachments/document/1 -d "{
  \"_name\"    : \"test.doc\",
  \"content\"  : \"$(openssl base64 -in test.doc)\"
}"

echo;
curl -X POST http://localhost:9200/test_attachments/_refresh

echo; echo ">>> Search for author 'John'"
curl "http://localhost:9200/test_attachments/_search?pretty=true&q=content.author:john&fields=content.title,content.author"

