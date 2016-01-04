package middlewares

import (
	"../models"
	"gopkg.in/olivere/elastic.v2"
)

// Create a new Index into Elasticsearch
func ImportTextElastic(inputText string, elasticURL string, logFile string) {
	log.SetLogger("file", logFile)
	log.Trace("Create a new Index in Elasticsearch")

	if inputText == "" {
		inputText = "test string"
	}

	log.Trace("inputText = %s ", inputText)

	// Obtain a client
	log.Trace("Create an Elasticsearch client : %s", elasticURL)
	client, err := elastic.NewClient(elastic.SetURL(elasticURL), elastic.SetSniff(false))
	if err != nil {
		log.Error("err : %s", err)
	}

	// Use the IndexExists service to check if a specified index exists.
	log.Trace("Search the postindex in Elasticsearch")
	exists, err := client.IndexExists("postindex").Do() // index should be in lower case
	if err != nil {
		log.Error("err : %s", err)
	}

	if !exists { // Create a new index
		createIndex, err := client.CreateIndex("postindex").Do()
		if err != nil {
			log.Error("err : %s", err)
		}
		if !createIndex.Acknowledged {
			log.Trace("Not ackowledged")
		}
	}

	post := models.Post{}
	post.UserId = 201
	//post.Content.TextMessage = inputText
	put1, err := client.Index().
		Index("postindex").
		Type("text").
		Id("Id").
		BodyJson(post).
		Do()

	if err != nil {
		log.Error("err : %s", err)
	}
	log.Trace("Indexed post %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index("postindex").Do()
	if err != nil {
		log.Error("err : %s", err)
	}
}
