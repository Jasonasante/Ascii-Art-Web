package main

// This simply declares what package this code is a part of. In our case, we are declaring it
// as part of the main package because we intend to have our application start by
// running the main() function in this file.

import (
	"fmt"
	"net/http"
	//This package is used to both create
	//an application capable of responding to web requests, as well as making web
	//requests to other servers.
)

// this function processes incoming web requests.
// Everytime someone visits the website, the code in this function gets run and determines
//what is returned to the viewer (there are other handlers too).
//all handlers have the same two elements
// 1) http.ResponseWriter
// allows us to modify the response that we want to send to whoever visited our website.

// 2)the pointer: *http.Request
// this  accesses data from the web request. For example, we might use this to get the users email address and
//password after they sign up for our web application.
func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Ascii Art GUI Generator!</h1>")

}

func main() {
	//this handles the web request to the server with the path /
	//so it covers all paths that the user may visit on the website and it would be processed by handlerFunc
	//for example: t http://localhost:3000/some-other-path
	http.HandleFunc("/", handlerFunc)
	// starts up a web server listening on port 8080 using the default http handlers
	// so when we run the file, we open a browser and type: "http://localhost:8080/"
	// which is saying 'try to load a web page from this computer at port 8080'
	http.ListenAndServe(":8080", nil)
}
