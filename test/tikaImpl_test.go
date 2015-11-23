package main

import (
    "bytes"

	"errors"
    "io"
    "strings"

	"os/exec"
    "gopkg.in/olivere/elastic.v2"
    "os"
    "time"
    
)


type Tweet struct {
  User     string    `json:"user"`
  Message  string    `json:"message"`
}

// use Apache Tika convert doc: ppt,pdf,word to plain text
func tikaImplt() {

    log.Trace("Convert document into plain text")
    
    // Create io.Reader 
    file, err := os.Open("./test.pdf")
    defer func() {
        file.Close()
    }()
    
    if err != nil {
        panic(err)
    }
    
    // Create io.Writer 
    ws := &bytes.Buffer{}
    DocToText(file, ws)
 
    // import a document to Elasticsearch 
    createIndexFromText(ws.String())
}

// Create a new Index in Elasticsearch
func createIndexFromText(post string) {

	log.Trace("Create a new Index in Elasticsearch")
	// Obtain a client
	client, err := elastic.NewClient(elastic.SetURL(elasticURL), elastic.SetSniff(false))
    log.Trace("Create a client to Elasticsearch")
	if err != nil {
		log.Error("err : %s", err)
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("twitter").Do() // index should be in lower case
    log.Trace("Search the attachment in Elasticsearch")
	if err != nil {
		log.Error("err : %s", err)
	}

	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("twitter").Do()
		if err != nil {
			log.Error("err : %s", err)
		}
		if !createIndex.Acknowledged {
			log.Trace("Not ackowledged")
		}
	}

    // Add a document to the index, using JSON serialization
    tweet := Tweet{User: "olivere", Message: post}
    _, err = client.Index().
        Index("twitter").
        Type("tweet").
        Id("1").
        BodyJson(tweet).
        Do()
    if err != nil {
        // Handle error
        panic(err)
    }

	if err != nil {
    	panic(err)
	}
}


// Convert document to plain text
func DocToText(in io.Reader, out io.Writer) error {
    log.SetLogger("file", logFileName)
    log.Trace("DocToText: Convert document to plain text")
    
	cmd := exec.Command("java", "-jar", "./bin/tika-app-1.7.jar", "-t")
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

func IsSupport(fileName string) bool {
	paths := strings.Split(fileName, "/")
	fileType := strings.Split(paths[len(paths)-1], ".")

	switch fileType[len(fileType)-1] {
	case "doc", "docx", "xls", "xlsx", "ppt", "pptx", "pdf", "epub", "html", "xml":
		return true
	}
	return false

}

func testReadWrite() {

	// open files r and w
	src, err := os.Open("/Users/huazhang/git/channel-service/code/input.txt")
	if err != nil {
		panic(err)
	}
	defer src.Close()

	des, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer des.Close()

	// do the actual work
	n, err := io.Copy(des, src) // <------ here !
	if err != nil {
		panic(err)
	}

	fmt.Printf("Copied %v bytes\n", n)

}


// Ref
// 1. http://rny.io/rails/elasticsearch/2013/08/05/full-text-search-for-attachments-with-rails-and-elasticsearch.html
// 2. Wrapper Apache Tika https://github.com/plimble/gika