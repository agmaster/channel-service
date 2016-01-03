// Go provides built-in support for [base64
// encoding/decoding](http://en.wikipedia.org/wiki/Base64).

package main

// This syntax imports the `encoding/base64` package with
// the `b64` name instead of the default `base64`. It'll
// save us some space below.
import (
    b64 "encoding/base64"
    curl "github.com/andelf/go-curl"
    "fmt"
    "time"
)

const POST_DATA = "a_test_data_only1"

// const POST_DATA ='{
//   "fields": [ "file.content_type" ],
//   "query": {
//     "match": {
//       "file.content_type": "text plain"
//     }
//   }
// }'
 

var sent = false


func sendrecv() {
    // ref : https://github.com/andelf/go-curl/blob/c965868dde67fef2abe524c33f1661d9ad233fac/examples/sendrecv.go
	easy := curl.EasyInit()
	defer easy.Cleanup()

	easy.Setopt(curl.OPT_URL, "http://localhost")

	easy.Setopt(curl.OPT_PORT, 9200)
	easy.Setopt(curl.OPT_VERBOSE, true)
	easy.Setopt(curl.OPT_CONNECT_ONLY, true)

	easy.Setopt(curl.OPT_WRITEFUNCTION, nil)

	if err := easy.Perform(); err != nil {
		println("ERROR: ", err.Error())
	}

	easy.Send([]byte("HEAD / HTTP/1.0\r\nHost: localhost\r\n\r\n"))

	buf := make([]byte, 1000)
	time.Sleep(1000000000) // wait gorotine
	num, err := easy.Recv(buf)
	if err != nil {
		println("ERROR:", err.Error())
	}
	println("recv num = ", num)
	// NOTE: must use buf[:num]
	println(string(buf[:num]))

	fmt.Printf("got:\n%#v\n", string(buf[:num]))
}

func main() {
    sendrecv()
}

func postCallback() {
    
    easy := curl.EasyInit()
    defer easy.Cleanup()
    
    //easy.Setopt(curl.OPT_URL, "http://localhost:9200/blog/post/3?pretty=true")

    easy.Setopt(curl.OPT_URL, "http://localhost")
    easy.Setopt(curl.OPT_PORT, 9200)
    easy.Setopt(curl.OPT_PUT, true)
    	easy.Setopt(curl.OPT_VERBOSE, true)
        
        easy.Setopt(curl.OPT_READFUNCTION,
        		func(ptr []byte, userdata interface{}) int {
        			// WARNING: never use append()
        			if !sent {
        				sent = true
        				ret := copy(ptr, POST_DATA)
        				return ret
        			}
        			return 0 // sent ok
        		})

        	// disable HTTP/1.1 Expect 100
        	easy.Setopt(curl.OPT_HTTPHEADER, []string{"Expect:"})
        	// must set
        	easy.Setopt(curl.OPT_POSTFIELDSIZE, len(POST_DATA))

        	if err := easy.Perform(); err != nil {
        		println("ERROR: ", err.Error())
        	}

        	time.Sleep(1000) // wait gorotine
    
    // // make a callback function
    easy.Setopt(curl.OPT_URL, "http://localhost:9200/test/1?pretty=true")
    fooTest := func (buf []byte, userdata interface{}) bool {
        println("DEBUG: size=>", len(buf))
        println("DEBUG: content=>", string(buf))
        return true
    }

    easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)
    if err := easy.Perform(); err != nil {
        fmt.Printf("ERROR: %v\n", err)
    }
    

    // Here's the `string` we'll encode/decode.
    data := "abc123!?$*&()'-=@~"

    // Go supports both standard and URL-compatible
    // base64. Here's how to encode using the standard
    // encoder. The encoder requires a `[]byte` so we
    // cast our `string` to that type.
    sEnc := b64.StdEncoding.EncodeToString([]byte(data))
    fmt.Println(sEnc)

    // Decoding may return an error, which you can check
    // if you don't already know the input to be
    // well-formed.
    sDec, _ := b64.StdEncoding.DecodeString(sEnc)
    fmt.Println(string(sDec))
    fmt.Println()

    // This encodes/decodes using a URL-compatible base64
    // format.
    uEnc := b64.URLEncoding.EncodeToString([]byte(data))
    fmt.Println(uEnc)
    uDec, _ := b64.URLEncoding.DecodeString(uEnc)
    fmt.Println(string(uDec))
}