package handlers

import (
	"encoding/json"
	"fmt"
    "log"
	"net/http"
    "os"
    "reflect"

	"../models"
    //"../elastics"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/julienschmidt/httprouter"
    "gopkg.in/olivere/elastic.v2"
    
)


type (
	// ProductController represents the controller for operating on the Product resource
	ProductController struct {
		session *mgo.Session
	}
)

// NewProductController provides a reference to a ProductController with provided mongo session
func NewProductController(s *mgo.Session) *ProductController {
	return &ProductController{s}
}

// Create a new Index in Elasticsearch 
func CreateIndex(product models.Product) {
	errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)

	// Obtain a client. You can provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetErrorLog(errorlog),  elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}
    
    // Use the IndexExists service to check if a specified index exists.
    exists, err := client.IndexExists("productindex").Do()
    if err != nil {
        panic(err)
    }
    
    if !exists {
        // Create a new index.
        createIndex, err := client.CreateIndex("productIndex").Do()
        if err != nil {
            panic(err)
        }
        if !createIndex.Acknowledged {
            fmt.Println("Not ackowledged")
        }
    }
    
    // Index a product (using JSON serialization) 
  //  product1 := models.Product{Name : "peanuts", Description: "good food for afternoon"}
    put1, err := client.Index().
        Index("productindex").
        Type("text").
        Id("1").
        BodyJson(product).
        Do()
    
        if err != nil {
            panic(err)
        }
        fmt.Printf("Indexed product %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
        
        
    	// Flush to make sure the documents got written.
    	_, err = client.Flush().Index("productindex").Do()
    	if err != nil {
    		panic(err)
    	}
    
}

func SearchIndexWithId(id string) {
	errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)

	// Obtain a client. You can provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetErrorLog(errorlog),  elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}
    
    get1, err := client.Get().
        Index("productindex").
       // Type("text").
        Id(id).
        Do()
    
    if err != nil {
        panic(err)
    }
    if get1.Found {
        //product = get1.(models.Product)
        fmt.Printf("Got document %s in verion %d from index %s \n", get1.Id, get1.Version, get1.Index )
          
    }


}

// Search with a term query in Elasticsearch
func SearchIndexWithTermQuery()(product  models.Product) {
	errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)

	// Obtain a client. You can provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetErrorLog(errorlog),  elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}
  
    // q := ""
    // if (q != "") {  //  GET:   /v1/posts[?limit=xx&offset=xx&q=xx]    q is a search string
    //      termQuery = elastic.NewTermQuery("name", q)
    // }
    //
    
    
    
    termQuery := elastic.NewTermQuery("name", "boil")
    	searchResult, err := client.Search().
    		Index("productindex").   // search in index "productindex"
    		Query(&termQuery).  // specify the query
    		Sort("name", true). // sort by "name" field, ascending
    		From(0).Size(10).   // take documents 0-9
    		Pretty(true).       // pretty print request and response JSON
    		Do()                // execute
    	if err != nil {
    		panic(err)
    	}
    
    
    
    // searchResult is of type SearchResult and returns hits, suggestions,
    // and all kinds of other information from Elasticsearch.
    fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	var ttyp models.Product 
   if searchResult.Hits != nil {
    	// TotalHits is another convenience function that works even when something goes wrong.
        // fmt.Printf("Found a total of %d products\n", searchResult.TotalHits())
    	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
    		t := item.(models.Product)
        
    		fmt.Printf("Product Name:  %s,  Description: %s, Link: %s\n", t.Name, t.Description, t.Permalink)
            product = t    
    	}
    }  else {
        // Not hits
        fmt.Print("Found no products\n")
    }

    return product
}

// GetProduct retrieves an individual post resource
func (uc ProductController) GetProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    fmt.Print("Get Product\n")
	// Grab id
     id := p.ByName("id")
     fmt.Println(id)

  
     // q is a search string
     //q := p.ByName("q")
    
    // // Verify id is ObjectId, otherwise bail
    // if !bson.IsObjectIdHex(id) {
    //     w.WriteHeader(404)
    //     return
    // }


	// Grab id
    //	oid := bson.ObjectIdHex(id)

	// Stub product
	u := models.Product{}

	// Fetch product from MongoDB
    /*
	if err := uc.session.DB("post_message_service").C("products").FindId(oid).One(&u); err != nil {
		w.WriteHeader(404)
		return
	} */
    
     // Fetch product from Elasticsearch
    // if id != "" {
 //        u = SearchIndexWithId(id);
 //    } else {
 //        u = SearchIndexWithTermQuery()
 //    }
      u = SearchIndexWithTermQuery()
    fmt.Println(u)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)


	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// CreateProduct creates a new post resource
func (uc ProductController) CreateProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an post to be populated from the body
	u := models.Product{}

	// Populate the post data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	//u.Id = bson.NewObjectId()
    u.Id =  1 //p.ByName("id")

	// Write the post to mongo
	uc.session.DB("post_message_service").C("products").Insert(u)
    
    // Write the product to Elasticsearch
     //product1 := models.Product{Name : "peanuts", Description: "good food for afternoon", Permalink: "www.google.com"}
     fmt.Printf("\nInsert Product name : %s , Description: %s, link: %s\n", u.Name, u.Description, u.Permalink)
     CreateIndex(u) 

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

// RemoveProduct removes an existing post resource
func (uc ProductController) RemoveProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
	if err := uc.session.DB("post_message_service").C("products").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}