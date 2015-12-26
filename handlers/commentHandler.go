package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	Logger "github.com/astaxie/beego/logs"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"../models"
)

var log = Logger.NewLogger(10000)

// NewCommentController provides a reference to a Controller with provided mongo session
func NewCommentController(s *mgo.Session, config Configuration, logFile string) *Controller {
	return &Controller{s, config, logFile}
}

// GetComment retrieves an individual comment resource
// handler.GetComment
func (uc Controller) GetComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.SetLogger("file", uc.logFile)
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	log.Trace("GetCommentWithQuery: retrieves an individual comment resource")

	u := models.Comment{}

	// Fetch comment
	if err := uc.session.DB(uc.config.Database).C("comments").FindId(oid).One(&u); err != nil {
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
func (uc Controller) CreateComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.SetLogger("file", uc.logFile)
	log.Trace("creates a new comment for a post")
	// Stub an comment to be populated from the body
	u := models.Comment{}

	// Populate the comment data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.Id = bson.NewObjectId()

	// Write the comment to mongo
	uc.session.DB(uc.config.Database).C("comments").Insert(u)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

// RemoveComment removes an existing comment resource
func (uc Controller) RemoveComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.SetLogger("file", uc.logFile)
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
	if err := uc.session.DB(uc.config.Database).C("comments").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}
