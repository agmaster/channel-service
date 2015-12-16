package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"
    
    "../models"
    "../middlewares"
)

var elasticURL = "http://127.0.0.1:9200"
var dbName = "channel_service"

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

// CreatePost creates a new post resource
func (uc PostController) CreatePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.SetLogger("file", logFileName)
	log.Trace("Create Post")

	// Specify the Mongodb database
	db := uc.session.DB(dbName)

	// Stub an post to be populated from the body
	u := models.Post{}

	// Populate the post data
	fmt.Print(r.Body)
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.Id = bson.NewObjectId()
	u.Active = true

	// Write the post to MongoDB collection-posts
	db.C("posts").Insert(u)

	log.Trace("Insert Post user-id : %d , type: %s, active : %t", u.UserId, u.Type, u.Active)
	CreateIndex(u)

	// store file, video, audio, and images into MongoDB
	if u.Type != "text" {
		saveFileToMongo(u, uc)
	}

	if u.Type == "file" {
		// Index files("doc", "docx", "xls", "xlsx", "ppt", "pptx", "pdf", "epub") to Elasticsearch
		log.Trace(" Index files into Elasticsearch, method : %s", r.Method)

		inputFile, handler, err := r.FormFile("filename")
		if err != nil {
			log.Error("err : %s", err)
		}
		log.Trace("handler.Header %s", handler.Header)

		// Create io.Writer
		outText := &bytes.Buffer{}

		middlewares.DocToText(inputFile, outText)
		importTextElastic(outText.String())
	}

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)
	fmt.Fprintf(w, "%s", uj)
}

// Store file into MongoDB via mgo
func saveFileToMongo(u models.Post, uc PostController) {
	log.SetLogger("file", logFileName)
	log.Trace(" store file into MongoDB")

	// Specify the Mongodb database
	db := uc.session.DB(dbName)

	// Capture multipart form file information
	/*
		file, handler, err := r.FormFile("filename")
		if err != nil {
			log.Error("err : %s", err)
		}
		log.Trace("handler.Header %s", handler.Header) */

	link := u.Link
	log.Trace("link = %s", link)

	//data, err := ioutil.ReadAll(link)
	data, err := ioutil.ReadFile(link)
	if err != nil {
		log.Error("err : %s", err)
	}

	// Create the file in the Mongodb Gridfs instance
	timestamp := time.Now().Unix()
	tmstr := strconv.FormatInt(timestamp, 10)
	fileName := "post-file-" + tmstr
	my_file, err := db.GridFS("posts").Create(fileName)

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

	//Write a log type message
	log.Trace("%d bytes written to the Mongodb instance", n)
}

// Create a new Index in Elasticsearch
func CreateIndex(post models.Post) {
	log.SetLogger("file", logFileName)
	log.Trace("Create a new Index in Elasticsearch %s", elasticURL)
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

	if !exists { // Create a new index.
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
	log.Trace("Indexed post %s to index %s, type %s in Elasticsearch", put1.Id, put1.Index, put1.Type)

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index("postindex").Do()
	if err != nil {
		log.Error("err : %s", err)
	}

}

// GetPost retrieves an individual post resource
// handler.GetPostWithQuery
func (uc PostController) GetPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
	log.Trace("Remove an existing post")
	// Grab id
	id := p.ByName("id")
	log.Trace("id = " + id)

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove post
	if err := uc.session.DB(dbName).C("posts").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(200) // Write status
}

// Search with a term query in Elasticsearch
func SearchIndexWithTermQuery(limit string, offset string, q string) (post models.Post) {
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
	var ttyp models.Post
	if searchResult.Hits != nil {
		// TotalHits is another convenience function that works even when something goes wrong.
		// log.Trace("Found a total of %d posts\n", searchResult.TotalHits())
		for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
			t := item.(models.Post)

			log.Trace("Post Name:  %s,  Type: %s, Active: %s\n", t.UserId, t.Type, t.Active)
			post = t
		}
	} else {
		// Not hits
		log.Trace("Found no posts\n")
	}

	return post
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

// Ref:  http://stackoverflow.com/questions/22159665/store-uploaded-file-in-mongodb-gridfs-using-mgo-without-saving-to-memory
