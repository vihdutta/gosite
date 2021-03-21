package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/projects", projects)
	fmt.Println(http.ListenAndServe(":6969", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("templates/index.html"))

	if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func projects(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("templates/projects.html"))

	if err := templates.ExecuteTemplate(w, "projects.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
