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
	"html/template"
	"os"
	"strings"
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
type Text struct {
	Input  string
	Output string
}

func homePage(w http.ResponseWriter, r *http.Request) {
	//if the handle is not / then it will return an error message
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>HTTP Status 404: Page Not Found</h1>")
		return
	}
	//t is the template file. ParseFiles opens up the template file and attempt to validate it.
	//If everything is correct there will be a nil error and a *template
	t, err := template.ParseFiles("template.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "<h1>HTTP Status 500: Internal Server Error(</h1>")
		fmt.Fprint(w, "<p>No banner Selected</p>")
		return
	}
	//t.Executes runs the data input (in this case Text{}) throught the template
	// which is then displayed on website via ResponseWriter
	err = t.Execute(w, Text{})
}

func asciiPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ascii-art" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>HTTP Status 404: Page Not Found</h1>")
		return
	}
	r.ParseForm()
	input := r.Form["input"][0]
	if input == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<h1>HTTP Status 400: Bad Request</h1>")
		fmt.Fprint(w, "<p>Empty input string</p>")
		return
	}
	for _, ele := range input {
		if (ele != 13) && (ele != 10) && (ele < 32 || ele > 126) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "<h1>HTTP Status 400: Bad Request</h1>")
			fmt.Fprint(w, "<p>Incorrect character detected</p>")
			return
		}
	}
	if r.FormValue("banner") == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<h1>HTTP Status 400: Bad Request</h1>")
		fmt.Fprint(w, "<p>No banner Selected</p>")
		return
	}
	banner := r.Form["banner"][0]
	_, err := os.ReadFile(banner + ".txt")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "<h1>HTTP Status 500: Internal Server Error(</h1>")
		fmt.Fprint(w, "<p>No banner file found</p>")
		return
	}
	var output string
	output = asciiArt(input, banner)
	p := Text{Input: input, Output: output}
	t, err := template.ParseFiles("result.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "<h1>HTTP Status 500: Internal Server Error</h1>")
		fmt.Fprint(w, "<p>No template found</p>")
		return
	}
	t.Execute(w, p)
	//t.Execute(w, Text{input, output}) this is another way of executing
}

func asciiArt(s string, b string) string {
	var emptyString string
	var inputString []string
	Content, _ := os.ReadFile(b + ".txt")
	asciiSlice2 := make([][]string, 95)
	s = strings.Replace(s, "\r\n", "\\n", -1)
	inputString = strings.Split(s, "\\n")
	for i := 0; i < len(asciiSlice2); i++ {
		asciiSlice2[i] = make([]string, 9)
	}
	var bubbleCount int
	count := 0
	for i := 1; i < len(Content); i++ {
		if Content[i] == '\n' && bubbleCount <= 94 {
			asciiSlice2[bubbleCount][count] = emptyString
			emptyString = ""
			count++
		}
		if count == 9 {
			count = 0
			bubbleCount++
		} else {
			if Content[i] != '\n' && Content[i] != '\r' {
				emptyString += string(Content[i])
			}
		}
	}
	var outputStr string
	var tempOutput [][]string
	for _, str := range inputString {
		for _, aRune := range str {
			tempOutput = append(tempOutput, asciiSlice2[aRune-rune(32)])
		}
		for i := range tempOutput[0] {
			for _, char := range tempOutput {
				outputStr += char[i]
			}
			outputStr += "\n"
		}
		tempOutput = nil
	}
	return outputStr
}

func main() {
	fmt.Println("Starting Server at Port 8080")
	fmt.Println("now open a broswer and enter: localhost:8080 into the URL")
	//this handles the web request to the server with the path /
	//so it covers all paths that the user may visit on the website and it would be processed by handlerFunc
	//for example: t http://localhost:3000/some-other-path
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ascii-art", asciiPage)
	// starts up a web server listening on port 8080 using the default http handlers
	// so when we run the file, we open a browser and type: "http://localhost:8080/"
	// which is saying 'try to load a web page from this computer at port 8080'
	http.ListenAndServe(":8080", nil)

}
