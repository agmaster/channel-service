package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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

// GetPost retrieves an individual post resource
func (uc PostController) GetPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Stub post
	u := Post{}

	// Fetch post
	if err := uc.session.DB("post_message_service").C("posts").FindId(oid).One(&u); err != nil {
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

// CreatePost creates a new post resource
func (uc PostController) CreatePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an post to be populated from the body
	u := Post{}

	// Populate the post data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.Id = bson.NewObjectId()

	// Write the post to mongo
	uc.session.DB("post_message_service").C("posts").Insert(u)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

// RemovePost removes an existing post resource
func (uc PostController) RemovePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
	if err := uc.session.DB("post_message_service").C("posts").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}
