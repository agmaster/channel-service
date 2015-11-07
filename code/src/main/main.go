package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

func main() {

	// Instantiate a new router
	router := httprouter.New()

	// Get a PostController instance
	handler := NewPostController(getSession())

	// Get total count of the posts
	router.GET("/v1/posts/count", handler.GetPostCount)

	// Get a post resource with query string
	// GET:   /v1/posts[?limit=xx&offset=xx&q=xx]    q is a search string
	router.GET("/v1/posts[?limit=xx&offset=xx&q=xx]", handler.GetPostWithQuery)
	// Create a new post
	router.POST("/v1/posts", handler.CreatePost)

	// Remove an existing post
	router.DELETE("/v1/posts/:id", handler.RemovePost)

	// Fire up the server
	http.ListenAndServe("127.0.0.1:3000", router)

}

// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
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
*/

/*  Insert a new product
 curl -XPOST -H 'Content-Type: application/json' -d '{"Name": "peanuts", "Description": "Honey Roasted peanuts", "MetaDescription": "Good taste food"}' http://127.0.0.1:3000/v1/products
Result : {"id":"563ab9e7d126195066000001","name":"peanuts","description":"Honey Roasted peanuts","permalink":"",
            "tax_category_id":0,"shipping_category_id":0,"deleted_at":"0001-01-01T00:00:00Z","meta_description":"",
            "meta_keywords":"","position":0,"is_featured":false,"can_discount":false,"distributor_only_membership":false}

 Query an existing product
curl -H "Content-Type: application/json" -X GET -v http://127.0.0.1:3000/v1/products/563ab9e7d126195066000001

*/
