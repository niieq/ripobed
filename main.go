package main

import (
	"html/template"

	"github.com/gorilla/mux"
	"net/http"
)

func mainRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	return r
}

func main() {

	r := mainRouter()
	http.ListenAndServe(":8080", r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, nil)

}
