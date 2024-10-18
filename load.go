package main

import (
	"fmt"
	"net/http"
	"text/template"
)

const port = ":8080"


func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	fmt.Println(" (http://localhost:8080) - Server Start on port :", port)
	http.ListenAndServe(port, nil)
}
