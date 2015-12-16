package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Post struct {
	Id     bson.ObjectId `json:"id" bson:"_id"`
	UserId int64         `json:"user-id"`
	Type   string        `json:"type"`
	Active bool          `json:"active"`
	Content
	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
}

type Content struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	TextMessage string `json:"text-message"`
}

type Posts []Post

/*
file content:   { "title": "xx",    "link": "xx",   "name": "xx",   "comment": "xx"}
text content:   { "text-message": "xx"}
video content:  { "title": "xx",    "link": "xx",   "name": "xx",   "comment": "xx"}
image content:  { "title": "xx",    "link": "xx",   "name": "xx",   "comment": "xx"}
*/

/*
#Insert a new post
 curl -XPOST -H 'Content-Type: application/json' -d '{"user-id": 101, "type": "text","active": true,  "text-message" : "Honey Roasted Peanuts" }' http://127.0.0.1:3000/v1/posts


 #upload a new image post
 curl -XPOST -H 'Content-Type: application/json' -d '{"user-id": 201, "type": "image","active": true,  "title" : "mylogo",  "comment" : "This is an image file" , "link" : "image=@/Users/huazhang/git/post-message-service/test/mylogo.jpg"}' http://127.0.0.1:3000/v1/posts


 #Query an existing post
curl -H "Content-Type: application/json" -X GET -v http://127.0.0.1:3000/v1/posts/563cf648d12619398ed7c71a
*/
