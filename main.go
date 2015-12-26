package main

import (
	"./handlers"
	"encoding/json"
	"fmt"
	Logger "github.com/astaxie/beego/logs"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
	"strings"
)

var log = Logger.NewLogger(10000)
var logFile = `{"filename":"channel-service.log"}`
var configFile = "./config.json"
var elasticURL = "http://127.0.0.1:9200"
var server = "http://127.0.0.1:3000"
var dbName = "channel_service"

func readConf(configFile string) (configuration handlers.Configuration, err error) {

	file, _ := os.Open(configFile)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return
	}

	return
}

func main() {

	//`os.Args[1:]`holds the arguments to the program.
	var s, sep string
	for i := 0; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)

	if len(os.Args) > 1 && os.Args[1] != "" {
		configFile = os.Args[1]
		fmt.Printf("configFile = %s \n", configFile)
	}

	if len(os.Args) > 2 && os.Args[2] != "" {
		// func Replace(s, old, new string, n int) string
		logFile = strings.Replace(logFile, "channel-service.log", os.Args[2], 1)
		fmt.Printf("logFile = %s \n", logFile)
	}

	// Initialize log variable (10000 is the cache size)
	log := Logger.NewLogger(10000)
	//log.SetLogger("console", `{"level":1}`)
	log.SetLogger("file", logFile)

	// read elasticsearch, mongodb setting from config.json
	configuration, err := readConf(configFile)

	// Instantiate a new router
	router := httprouter.New()

	// Get a PostController instance
	//uc := getController(configuration, logFile)

	postHandler := handlers.NewPostController(getSession(), configuration, logFile)

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
	router.GET("/v1/posts/count", postHandler.GetPostCount)

	// Get a post resource with query string
	// GET:   /v1/posts[?limit=xx&offset=xx&q=xx]    q is a search string
	router.GET("/v1/posts", postHandler.GetPost)
	// Create a new postname := value
	router.POST("/v1/posts", postHandler.CreatePost)

	// Remove an existing post
	router.PUT("/v1/posts/:id", postHandler.RemovePost)

	// Get a CommentController instance
	commentHandler := handlers.NewCommentController(getSession(), configuration, logFile)

	// Create a new comment for a post
	router.POST("/v1/posts/:post-id/comments", commentHandler.CreateComment)

	// Get :/v1/posts/:post-d/comments
	//router.GET("/v1/posts/:post-id/comments", commentHandler.GetComment)

	// Get a UploadFileHandler instance
	//	uploadFileController := handlers.NewUploadFileController(getSession(dbName))
	//	router.POST("/v1/uploadfile", uploadFileController.UploadFile)

	// Fire up channel-serivce IP and port
	log.Trace("start web service on %s \n", configuration.Server)
	http.ListenAndServe(configuration.Server, router)
}

// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {
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
