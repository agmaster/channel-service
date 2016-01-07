package handlers

import (
	"../middlewares"
	"../models"
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
	"os"
	"reflect"
	"strconv"
	"time"

	"errors"
	"io"
	"os/exec"
)

// Controller represents the controller for operating on the Post resource
type Controller struct {
	session *mgo.Session
	config  Configuration
	logFile string
}

type Configuration struct {
	Elasticsearch string
	Database      string
	Server        string
}

// NewPostController provides a reference to a Controller with provided mongo session
func NewPostController(s *mgo.Session, config Configuration, logFile string) *Controller {

	return &Controller{s, config, logFile}
}

// CreatePost creates a new post resource
func (uc Controller) CreatePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.SetLogger("file", uc.logFile)
	log.Trace("Create Post")

	// Specify the Mongodb database
	db := uc.session.DB(uc.config.Database)

	// Stub an post to be populated from the body
	u := models.Post{}

	// Populate the post data
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Error("decode request body failed")
	}

	// Add an Id
	u.Id = bson.NewObjectId()
	u.Active = true
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	fmt.Println(u)
	log.Debug("u = %v", u)

	// Write the post to MongoDB collection-posts

	if err := db.C("posts").Insert(u); err != nil {
		log.Error("write the post to MongoDB failed. Error : %v", err)
		return
	}

	log.Debug("Insert Post user-id : %d , type: %s, active : %t", u.UserId, u.Type, u.Active)

	// Create a new Index in Elasticsearch
	CreateIndex(u, uc)

	// store file, video, audio, and images into MongoDB
	if u.Type == "file" || u.Type == "video" || u.Type == "image" {

		saveFileToMongo(u, uc)

	}

	// Index files to Elasticsearch
	if u.Type == "file" { //to be updated  2016.1.2
		log.Debug(" Index files into Elasticsearch, method : %s", r.Method)
		log.Trace("u.Content.Link = %s ", u.Content.Link)

		file := "/Users/huazhang/git/channel-service/test/test.txt"

		f, err := os.Open(file)
		if err != nil {
			//log.Error("err : %v", err)
			panic(err)
		}

		// Create io.Writer
		ws := &bytes.Buffer{}

		//DocToText(f, ws)
		middlewares.DocToText(f, ws, uc.logFile)
		log.Trace("ws.String() = %s", ws.String())
		fmt.Println(ws.String())
		// assert.NoError(t, err)
		// assert.True(t, ws.Len() > 0)
		middlewares.ImportTextElastic(ws.String(), uc.config.Elasticsearch, uc.logFile)

		/*
			inputFile, err := os.Open(u.Content.Link)
			//	inputFile, handler, err := r.FormFile("filename")
			if err != nil {
				log.Error("err : %s", err)
			}

			// log.debug("handler.Header %s", handler.Header)

			// Create io.Writer
			outText := &bytes.Buffer{}

			// use tika to convert "doc", "docx", "xls", "xlsx", "ppt", "pptx", "pdf", "epub" to text
			//middlewares.DocToText(inputFile, outText, uc.logFile)
			DocToText(inputFile, outText)

			// index text into Elasticsearch
			middlewares.ImportTextElastic(outText.String(), uc.config.Elasticsearch, uc.logFile)*/

	}

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)
	log.Debug("Store the post into MongoDB and Elasticsearch")
	log.Debug("post : %s", uj)
	fmt.Fprintf(w, "%s", uj)

	//log.Trace(w, "%s", uj)
}

// Convert document to plain text
func DocToText(in io.Reader, out io.Writer) error {

	cmd := exec.Command("java", "-jar", "./lib/tika-app-1.7.jar", "-t")
	stderr := bytes.NewBuffer(nil)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = stderr

	cmd.Start()
	cmdDone := make(chan error, 1)
	go func() {
		cmdDone <- cmd.Wait()
	}()

	select {
	case <-time.After(time.Duration(500000) * time.Millisecond):
		if err := cmd.Process.Kill(); err != nil {
			return errors.New(err.Error())
		}
		<-cmdDone
		return errors.New("Command timed out")
	case err := <-cmdDone:
		if err != nil {
			return errors.New(stderr.String())
		}
	}

	return nil
}

// Store file, video, audio, and images into MongoDB via mgo
func saveFileToMongo(u models.Post, uc Controller) {
	log.SetLogger("file", uc.logFile)
	log.Trace("store file into MongoDB : %s ", uc.config.Database)

	// Specify the Mongodb database
	db := uc.session.DB(uc.config.Database)

	// Capture multipart form file information
	/*
		file, handler, err := r.FormFile("filename")
		if err != nil {
			log.Error("err : %s", err)
		}
		log.Trace("handler.Header %s", handler.Header) */
	log.Debug("u.Content.Link = %s", u.Content.Link)
	link := u.Content.Link

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
func CreateIndex(post models.Post, uc Controller) {
	log.SetLogger("file", uc.logFile)
	log.Trace("Create a new Index in Elasticsearch %s", uc.config.Elasticsearch)
	// Obtain a client
	client, err := elastic.NewClient(elastic.SetURL(uc.config.Elasticsearch), elastic.SetSniff(false))
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
func (uc Controller) GetPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.SetLogger("file", uc.logFile)
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
func (uc Controller) GetPostCount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	log.SetLogger("file", uc.logFile)
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
func (uc Controller) RemovePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.SetLogger("file", uc.logFile)
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
	if err := uc.session.DB(uc.config.Database).C("posts").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(200) // Write status
}

// Search with a term query in Elasticsearch
func SearchIndexWithTermQuery(limit string, offset string, q string, uc Controller) (post models.Post) {
	log.SetLogger("file", uc.logFile)
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

func SearchIndexWithId(id string, uc Controller) {
	log.SetLogger("file", uc.logFile)

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
