package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/", Root)
	fmt.Println("El middleware esta en linea en localhost:8000")
	http.ListenAndServe(":8000", nil)
}

func LoadHTML(filename string) string {
	html, _ := ioutil.ReadFile(filename)
	return string(html)
}

func Root(res http.ResponseWriter, _ *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		LoadHTML("../static/middleware/index.html"),
	)
}
