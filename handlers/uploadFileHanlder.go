package handlers

import (
	"bytes"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/olivere/elastic.v2"
	"io/ioutil"
	"net/http"

      "../models"
      "../middlewares"
)

type (
	// UploadFileController represents the controller for uploading files
	UploadFileController struct {
		session *mgo.Session
	}
)

// NewUploadFileController provides a reference to a UploadFileController with provided mongo session
func NewUploadFileController(s *mgo.Session) *UploadFileController {
	return &UploadFileController{s}
}

/* Ref: The second example in the documentation for GridFS.Create
curl -i -X POST -H "Content-Type: multipart/form-data" \
-F "filename=@/Users/huazhang/git/channel-service/test/test.pdf" -v http://127.0.0.1:3000/v1/uploadfile
*/
func (uc UploadFileController) UploadFile(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.SetLogger("file", logFileName)
	log.Trace(" Upload file into mongod, method : %s", r.Method)

	inputFile, handler, err := r.FormFile("filename")
	if err != nil {
		log.Error("err : %s", err)
	}
	log.Trace("handler.Header %s", handler.Header)

	// Create io.Writer
	outText := &bytes.Buffer{}

	middlewares.DocToText(inputFile, outText)
	importTextElastic(outText.String())
	saveFile2Mongodb(uc, w, r)
}

// Create a new Index into Elasticsearch
func importTextElastic(inputText string) {
	log.SetLogger("file", logFileName)
	log.Trace("Create a new Index in Elasticsearch")

	// Obtain a client
	log.Trace("Create an Elasticsearch client")
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
	post.TextMessage = inputText
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

// Store file into MongoDB via mgo
func saveFile2Mongodb(uc UploadFileController, w http.ResponseWriter, r *http.Request) {
	log.SetLogger("file", logFileName)
	log.Trace(" store file into MongoDB")

	// Specify the Mongodb database
	db := uc.session.DB(dbName)

	// Capture multipart form file information
	file, handler, err := r.FormFile("filename")
	if err != nil {
		fmt.Println(err)
	}
	log.Trace("handler.Header %s", handler.Header)

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error("err : %s", err)
	}

	// Create the file in the Mongodb Gridfs instance
	my_file, err := db.GridFS("posts").Create("post_10002")
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
	log.Trace("%d bytes written to the Mongodb instance\n", n)
}

// Ref:  http://stackoverflow.com/questions/22159665/store-uploaded-file-in-mongodb-gridfs-using-mgo-without-saving-to-memory
