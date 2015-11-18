package main

import (
	// "strings"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

// Ref: The second example in the documentation for GridFS.Create
// test command:   curl -i -X POST -H "Content-Type: multipart/form-data" \
-F "filename=@/Users/huazhang/test.txt" -v http://127.0.0.1:3000/v1/uploadfile

func (uc UploadFileController) UploadFile(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.SetLogger("file", logFileName)
	log.Trace(" Upload file into mongod, method : %s", r.Method)
	
    //mgoTestCode(uc)
    
	// Specify the Mongodb database
	db := uc.session.DB("channel_service")

	// Capture multipart form file information
	file, handler, err := r.FormFile("filename")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("handler.Header %s", handler.Header)

	// Read the file into memory
	data, err := ioutil.ReadAll(file)
	//data, err := ioutil.ReadFile(fileName)
	// ... check err value for nil
	if err != nil {
		log.Error("err : %s", err)
	}

	// Create the file in the Mongodb Gridfs instance
	my_file, err := db.GridFS("posts").Create("post_10001")
	// ... check err value for nil
	if err != nil {
		log.Error("err : %s", err)
	}

	// Write the file to the Mongodb Gridfs instance
	n, err := my_file.Write(data)
	// ... check err value for nil
	if err != nil {
		log.Error("err : %s", err)
	}

	// Close the file
	err = my_file.Close()
	// ... check err value for nil
	if err != nil {
		log.Error("err : %s", err)
	}

	// Write a log type message
	fmt.Printf("%d bytes written to the Mongodb instance\n", n)

	// ... other statements redirecting to rest of user flow...
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
// Ref: The second example in the documentation for GridFS.Create
// test command:  curl -i -F name=test -F filedata=@test.txt http://127.0.0.1/v1/uploadfile
func (uc UploadFileController) UploadFile2(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.SetLogger("file", logFileName)
	log.Trace(" Upload file into mongodb, method : %s", r.Method)
	
    //mgoTestCode(uc)
    
	// Specify the Mongodb database
	db := uc.session.DB("channel_service")

	// Capture multipart form file information
	file, handler, err := r.FormFile("filename")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("handler.Header %s", handler.Header)

	// Read the file into memory
	data, err := ioutil.ReadAll(file)
	//data, err := ioutil.ReadFile(fileName)
	// ... check err value for nil
	if err != nil {
		log.Error("err : %s", err)
	}

	// Create the file in the Mongodb Gridfs instance
	my_file, err := db.GridFS("posts").Create("post_10001")
	// ... check err value for nil
	if err != nil {
		log.Error("err : %s", err)
	}

	// Write the file to the Mongodb Gridfs instance
	n, err := my_file.Write(data)
	// ... check err value for nil
	if err != nil {
		log.Error("err : %s", err)
	}

	// Close the file
	err = my_file.Close()
	// ... check err value for nil
	if err != nil {
		log.Error("err : %s", err)
	}

	// Write a log type message
	fmt.Printf("%d bytes written to the Mongodb instance\n", n)

	// ... other statements redirecting to rest of user flow...
}

func mgoTestCode(uc UploadFileController) {
	log.SetLogger("file", logFileName)
	log.Trace(" creat a temp file for mgo Test")

	// Specify the Mongodb database
	db := uc.session.DB("channel_service")
	file, err := db.GridFS("fs").Create("myfile.txt")

	n, err := file.Write([]byte(" Mgo test done!"))
	if err != nil {
		log.Error("err : %s", err)
	}

	err = file.Close()
	if err != nil {
		log.Error("err : %s", err)
	}

	log.Trace("%d bytes written\n", n)

}

// Ref:  http://stackoverflow.com/questions/22159665/store-uploaded-file-in-mongodb-gridfs-using-mgo-without-saving-to-memory
