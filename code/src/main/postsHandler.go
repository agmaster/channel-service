package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

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

var elasticURL = "http://127.0.0.1:9200"

// NewPostController provides a reference to a PostController with provided mongo session
func NewPostController(s *mgo.Session) *PostController {
	return &PostController{s}
}

// Create a new Index in Elasticsearch
func CreateIndex(post Post) {
	log.Trace("Create a new Index in Elasticsearch")
	// Obtain a client
	client, err := elastic.NewClient(elastic.SetURL(elasticURL), elastic.SetSniff(false))
    log.Trace("Create a client to Elasticsearch")
	if err != nil {
		log.Error("err : %s", err)
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("postindex").Do() // index should be in lower case
    log.Trace("Search the postindex in Elasticsearch")
	if err != nil {
		log.Error("err : %s", err)
	}

	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("postindex").Do()
		if err != nil {
			log.Error("err : %s", err)
		}
		if !createIndex.Acknowledged {
			log.Trace("Not ackowledged")
		}
	}

	// Index a post (using JSON serialization)
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

// CreatePost creates a new post resource
func (uc PostController) CreatePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.SetLogger("file", logFileName)
	// Stub an post to be populated from the body
	u := Post{}

	// Specify the Mongodb database
	db := uc.session.DB("channel_service")

	// Populate the post data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.Id = bson.NewObjectId()

	// Write the post to mongo
	db.C("posts").Insert(u)

	// store the file(.xls, .pdf, .docx, .txt, .rtf) into MongoDB
	if u.Type == "type" {
		// Capture multipart form file information
		file, handler, err := r.FormFile("filename")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("handler.Header %s", handler.Header)

		// Read the file into memory
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Error("err : %s", err)
		}

		// Create the file in the Mongodb Gridfs instance
		my_file, err := db.GridFS("posts").Create("post_10001")

		if err != nil {
			log.Error("err : %s", err)
		}

		// Write the file to the Mongodb Gridfs instance
		n, err := my_file.Write(data)

		if err != nil {
			log.Error("err : %s", err)
		}

		// Close the file
		err = my_file.Close()

		if err != nil {
			log.Error("err : %s", err)
		}

		// Write a log type message
		fmt.Printf("%d bytes written to the Mongodb instance\n", n)
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write the post to Elasticsearch
	log.Trace("\nInsert Post user-id : %d , type: %s, active : %t\n", u.UserId, u.Type, u.Active)
	CreateIndex(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

func SearchIndexWithId(id string) {
	log.SetLogger("file", logFileName)

	// Obtain a client. You can provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		log.Error("err : %s", err)
	}

	get1, err := client.Get().
		Index("postindex").
		// Type("text").
		Id(id).
		Do()

	if err != nil {
		log.Error("err : %s", err)
	}
	if get1.Found {
		//post = get1.(Post)
		log.Trace("Got document %s in verion %d from index %s \n", get1.Id, get1.Version, get1.Index)
	}
}

// Search with a term query in Elasticsearch
func SearchIndexWithTermQuery(limit string, offset string, q string) (post Post) {
	log.SetLogger("file", logFileName)
	log.Trace("Search with a term query in Elasticsearch")

	// Obtain a client. You can provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		// Handle error
		log.Error("err : %s", err)
	}
	termQuery := elastic.NewTermQuery("user-id", 101)

	if q != "" { //  GET:   /v1/posts[?limit=xx&offset=xx&q=xx]    q is a search string
		log.Trace(" q = %s", q)
		termQuery = elastic.NewTermQuery("user-id", q)
	}

	searchResult, err := client.Search().
		Index("postindex").    // search in index "postindex"
		Query(&termQuery).     // specify the query
		Sort("user-id", true). // sort by "name" field, ascending
		From(0).Size(10).      // take documents 0-9
		Pretty(true).          // pretty print request and response JSON
		Do()                   // execute
	if err != nil {
		log.Error("err : %s", err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	log.Trace("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	var ttyp Post
	if searchResult.Hits != nil {
		// TotalHits is another convenience function that works even when something goes wrong.
		// log.Trace("Found a total of %d posts\n", searchResult.TotalHits())
		for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
			t := item.(Post)

			log.Trace("Post Name:  %s,  Type: %s, Active: %s\n", t.UserId, t.Type, t.Active)
			post = t
		}
	} else {
		// Not hits
		log.Trace("Found no posts\n")
	}

	return post
}

func convertQueryStr(q string) (key string, val string) {
	log.SetLogger("file", logFileName)
	queryStr := strings.SplitN("user-id=101", "=", 2)
	key = queryStr[0]
	val = queryStr[1]

	log.Trace("key=%s, val=%s", key, val)

	return key, val
}

// GetPost retrieves an individual post resource
// handler.GetPostWithQuery
func (uc PostController) GetPostWithQuery(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.SetLogger("file", logFileName)
	log.Trace("GetPostWithQuery: retrieves an individual post resource")
	testIndexName := "postindex"

	var (
		q      = "*"
		limit  = 10
		offset = 0
	)

	queryForm, err := url.ParseQuery(r.URL.RawQuery)

	if err == nil && len(queryForm["limit"]) > 0 {
		fmt.Fprintln(w, queryForm["limit"])
		limit, err := strconv.Atoi(queryForm["limit"][0])
		if err != nil {
			log.Error("err : %s", err)
		}
		log.Trace("limit = %d", limit)
	}

	if err == nil && len(queryForm["offset"]) > 0 {
		fmt.Fprintln(w, queryForm["offset"])
		offset, err := strconv.Atoi(queryForm["offset"][0])
		if err != nil {
			log.Error("err : %s", err)
		}
		log.Trace("offset = %d", offset)
	}

	if err == nil && len(queryForm["q"]) > 0 {
		fmt.Fprintln(w, queryForm["q"])
		q = queryForm["q"][0]
	}

	log.Trace("q = %s ", q)
	//termQuery := elastic.NewMatchQuery(key, val)
	queryStringQuery := elastic.NewQueryStringQuery(q)

	// Obtain a client. You can provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetSniff(false))
	searchResult, err := client.Search().
		Index(testIndexName).
		Query(&queryStringQuery).
		Sort("id", true).         // sort by "user" field, ascending
		From(offset).Size(limit). // take documents 0-9
		Pretty(true).             // pretty print request and response JSON
		//Suggester(ts).
		Do()
	if err != nil {
		log.Error("err : %s", err)
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(searchResult)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// Get total count of the posts
// GET("/v1/posts/count", handler.GetPostCount)
func (uc PostController) GetPostCount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	log.SetLogger("file", logFileName)
	log.Trace("Get total count of the posts")

	// Obtain a client
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		log.Error("err : %s", err)
	}

	// Count documents
	//count, err := client.Count("postindex").Type("text").Do()
	count, err := client.Count("postindex").Do()
	if err != nil {
		log.Error("err : %s", err)
	}

	if count == 0 {
		w.WriteHeader(404)
		return
	}
	log.Trace("Found a total of %d posts\n", count)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "Found a total of %d posts\n", count)

}

// RemovePost removes an existing post resource
func (uc PostController) RemovePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.SetLogger("file", logFileName)
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
	if err := uc.session.DB("channel_service").C("posts").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}

/*
func TestSearchSourceInnerHits(t *testing.T) {
	matchAllQ := NewMatchAllQuery()
	builder := NewSearchSource().Query(matchAllQ).
		InnerHit("comments", NewInnerHit().Type("comment").Query(NewMatchQuery("user", "olivere"))).
		InnerHit("views", NewInnerHit().Path("view"))
	data, err := json.Marshal(builder.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"inner_hits":{"comments":{"type":{"comment":{"query":{"match":{"user":{"query":"olivere"}}}}}},"views":{"path":{"view":{}}}},"query":{"match_all":{}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

*/
