package main

import (
	"./handlers"
	"encoding/json"
	Logger "github.com/astaxie/beego/logs"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
)

var log = Logger.NewLogger(10000)
var logFileName = `{"filename":"channel-service.log"}`
var elasticURL = "http://127.0.0.1:9200"
var dbName = "channel_service"

type Configuration struct {
	Elasticsearch string
	Database      string
	Server        string
}

func readConf(fileName string) (elastic, dbstr, server string, err error) {

	file, _ := os.Open(fileName)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		return
	}
	// Elasticsearch server IP and port
	elastic = configuration.Elasticsearch

	// Mongodb connection string
	dbstr = configuration.Database

	// Channel-serivce IP and port
	server = configuration.Server

	return

}

func main() {

	// Initialize log variable (10000 is the cache size)
	log := Logger.NewLogger(10000)
	//log.SetLogger("console", `{"level":1}`)
	log.SetLogger("file", logFileName)

	fileName := "./config.json"
	elastic, dbstr, server, err := readConf(fileName)
	if err != nil {
		log.Trace("Read config file failed! ", err)
	}

	log.Trace("server = %s", server)
	elasticURL = elastic
	dbName = dbstr

	// Instantiate a new router
	router := httprouter.New()

	// Get a PostController instance
	handler := handlers.NewPostController(getSession(dbstr))

	// first create a client
	client, err := statsd.NewClient("127.0.0.1:8125", "test-client")
	// handle any errors
	if err != nil {
		log.Trace("Start statsd client failed! ", err)
	}
	// make sure to clean up
	defer client.Close()

	// Send a stat
	client.Inc("stat1", 42, 1.0)

	// Get total count of the posts
	router.GET("/v1/posts/count", handler.GetPostCount)

	// Get a post resource with query string
	// GET:   /v1/posts[?limit=xx&offset=xx&q=xx]    q is a search string
	router.GET("/v1/posts", handler.GetPost)
	// Create a new postname := value
	router.POST("/v1/posts", handler.CreatePost)

	// Remove an existing post
	router.PUT("/v1/posts/:id", handler.RemovePost)

	// Get a CommentController instance
	commentHandler := handlers.NewCommentController(getSession(dbstr))

	// Create a new comment for a post
	router.POST("/v1/posts/:post-id/comments", commentHandler.CreateComment)

	// Get :/v1/posts/:post-d/comments
	//router.GET("/v1/posts/:post-id/comments", commentHandler.GetComment)

	// Get a UploadFileHandler instance
	uploadFileController := handlers.NewUploadFileController(getSession(dbstr))
	router.POST("/v1/uploadfile", uploadFileController.UploadFile)

	// Fire up the server
	log.Trace("start web service on %s \n", server)
	http.ListenAndServe(server, router)
}

// getSession creates a new mongo session and panics if connection error occurs
func getSession(dbstr string) *mgo.Session {
	// Connect to our local mongo
	// "Database": "mongodb://localhost",
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		log.Trace("connection error~!", err)

	}

	// Deliver session
	return s
}

// Reference: 1. http://stevenwhite.com/building-a-rest-service-with-golang-3/
//            2. https://github.com/swhite24/go-rest-tutorial/blob/master/server.go

// test cases:

/*   Insert a new post
 curl -XPOST -H 'Content-Type: application/json' -d '{"user-id": 101, "type": "text", "content": "Hello World!"}' http://127.0.0.1:3000/v1/posts
Result: {"id":"563aa288d1261946cb000001","type":"text","content":"Hello World!","user-id":101}

Query an existing post
curl -H "Content-Type: application/json" -X GET -v http://127.0.0.1:3000/v1/posts/563aa288d1261946cb000001

Insert a new product
curl -XPOST -H 'Content-Type: application/json' -d '{"Name": "peanuts", "Description": "Honey Roasted peanuts", "MetaDescription": "Good taste food"}' http://127.0.0.1:3000/v1/products
Result : {"id":"563ab9e7d126195066000001","name":"peanuts","description":"Honey Roasted peanuts","permalink":"",
            "tax_category_id":0,"shipping_category_id":0,"deleted_at":"0001-01-01T00:00:00Z","meta_description":"",
            "meta_keywords":"","position":0,"is_featured":false,"can_discount":false,"distributor_only_membership":false}

 Query an existing product
curl -H "Content-Type: application/json" -X GET -v http://127.0.0.1:3000/v1/products/563ab9e7d126195066000001

*/
