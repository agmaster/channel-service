package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	Logger "github.com/astaxie/beego/logs"
    
    "../models"
)

var log = Logger.NewLogger(10000)
var logFileName = `{"filename":"channel-service.log"}`


type (
	// CommentController represents the controller for operating on the Comment resource
	CommentController struct {
		session *mgo.Session
	}
)

// NewCommentController provides a reference to a CommentController with provided mongo session
func NewCommentController(s *mgo.Session) *CommentController {
	return &CommentController{s}
}

// GetComment retrieves an individual comment resource
// handler.GetComment
func (uc CommentController) GetComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
    
    log.SetLogger("file", logFileName)
    log.Trace("GetCommentWithQuery: retrieves an individual comment resource")
    
    u := models.Comment {}
   
    // Fetch comment
    if err := uc.session.DB("channel_service").C("comments").FindId(oid).One(&u); err != nil {
        w.WriteHeader(404)
        return
    }
    
	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// CreateComment creates a new comment resource
func (uc CommentController) CreateComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    log.SetLogger("file", logFileName)
    log.Trace("creates a new comment for a post")
	// Stub an comment to be populated from the body
	u := models.Comment{}

	// Populate the comment data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.Id = bson.NewObjectId()

	// Write the comment to mongo
	uc.session.DB("channel_service").C("comments").Insert(u)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

// RemoveComment removes an existing comment resource
func (uc CommentController) RemoveComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	// Remove comment
	if err := uc.session.DB("channel_service").C("comments").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}
