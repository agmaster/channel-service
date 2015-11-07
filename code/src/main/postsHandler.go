package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v2"
)

type (
	// PostController represents the controller for operating on the Post resource
	PostController struct {
		session *mgo.Session
	}
)

// NewPostController provides a reference to a PostController with provided mongo session
func NewPostController(s *mgo.Session) *PostController {
	return &PostController{s}
}

// Create a new Index in Elasticsearch
func CreateIndex(post Post) {
	errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)

	// Obtain a client. You can provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetURL("http://127.0.0.1:9200"), elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("postindex").Do() // index should be in lower case
	if err != nil {
		panic(err)
	}

	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("postindex").Do()
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			fmt.Println("Not ackowledged")
		}
	}

	// Index a post (using JSON serialization)
	put1, err := client.Index().
		Index("postindex").
		Type("text").
		Id("1").
		BodyJson(post).
		Do()

	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed post %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index("postindex").Do()
	if err != nil {
		panic(err)
	}

}

func SearchIndexWithId(id string) {
	errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)

	// Obtain a client. You can provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}

	get1, err := client.Get().
		Index("postindex").
		// Type("text").
		Id(id).
		Do()

	if err != nil {
		panic(err)
	}
	if get1.Found {
		//post = get1.(Post)
		fmt.Printf("Got document %s in verion %d from index %s \n", get1.Id, get1.Version, get1.Index)

	}

}

// Search with a term query in Elasticsearch
func SearchIndexWithTermQuery() (post Post) {
	errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)

	// Obtain a client. You can provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}

	// q := ""
	// if (q != "") {  //  GET:   /v1/posts[?limit=xx&offset=xx&q=xx]    q is a search string
	//      termQuery = elastic.NewTermQuery("name", q)
	// }
	//

	termQuery := elastic.NewTermQuery("user-id", 101)
	searchResult, err := client.Search().
		Index("postindex").    // search in index "postindex"
		Query(&termQuery).     // specify the query
		Sort("user-id", true). // sort by "name" field, ascending
		From(0).Size(10).      // take documents 0-9
		Pretty(true).          // pretty print request and response JSON
		Do()                   // execute
	if err != nil {
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	var ttyp Post
	if searchResult.Hits != nil {
		// TotalHits is another convenience function that works even when something goes wrong.
		// fmt.Printf("Found a total of %d posts\n", searchResult.TotalHits())
		for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
			t := item.(Post)

			fmt.Printf("Post Name:  %s,  Type: %s, Active: %s\n", t.UserId, t.Type, t.Active)
			post = t
		}
	} else {
		// Not hits
		fmt.Print("Found no posts\n")
	}

	return post
}

// GetPost retrieves an individual post resource
func (uc PostController) GetPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	//oid := bson.ObjectIdHex(id)

	// Stub post
	u := Post{}

	// Fetch post
	/*if err := uc.session.DB("post_message_service").C("posts").FindId(oid).One(&u); err != nil {
		w.WriteHeader(404)
		return
	}*/

	// Fetch post from Elasticsearch
	u = SearchIndexWithTermQuery()
	fmt.Println(u)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// Get total count of the posts
//router.GET("/v1/posts/count", handler.GetPostCount)
func (uc PostController) GetPostCount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)

	// Obtain a client
	client, err := elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	// Count documents
	//count, err := client.Count("postindex").Type("text").Do()
	count, err := client.Count("postindex").Do()
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		w.WriteHeader(404)
		return
	}

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "Found a total of %d posts\n", count)
}

// handler.GetPostWithQuery
func (uc PostController) GetPostWithQuery(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

// CreatePost creates a new post resource
func (uc PostController) CreatePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an post to be populated from the body
	u := Post{}

	// Populate the post data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.Id = bson.NewObjectId()

	// Write the post to mongo
	uc.session.DB("post_message_service").C("posts").Insert(u)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write the post to Elasticsearch
	//product1 := models.Product{Name : "peanuts", Description: "good food for afternoon", Permalink: "www.google.com"}
	fmt.Printf("\nInsert Post user-id : %d , type: %s, active : %t\n", u.UserId, u.Type, u.Active)
	CreateIndex(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

// RemovePost removes an existing post resource
func (uc PostController) RemovePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove post
	if err := uc.session.DB("post_message_service").C("posts").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}
