package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Comment struct {
	Id     bson.ObjectId `json:"id" bson:"_id"`
	UserId int64         `json:"user-id"`
    PostId int64         `json:"post-id"`
	Active bool          `json:"active"`
	Content string       `json:"content"`
	CreatedAt time.Time `json:"created-at"`
}


type Comments []Comment

/*

 #Query an existing post
curl -H "Content-Type: application/json" -X GET -v http://127.0.0.1:3000/v1/posts/1/comments
*/
