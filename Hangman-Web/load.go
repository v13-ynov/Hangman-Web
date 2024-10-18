package main

import (
	"fmt"
	"net/http"
	"text/template"
)

const port = ":8080"

func Home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Home Page")
}

func About(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "About Page")
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("./templates/" + tmpl + ".page.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	fmt.Println(" (http://localhost:8080) - Server Start on port : ", port)
	http.ListenAndServe(port, nil)
}
