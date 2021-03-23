package main

import (
	"fmt"
	"html/template"
	"net/http"

	quote "github.com/vihdutta/gowebsite/modules"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/projects", projects)
	http.HandleFunc("/statistics", statistics)
	fmt.Println(http.ListenAndServe(":6969", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("templates/index.html"))
	fmt.Println("home")

	if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func projects(w http.ResponseWriter, r *http.Request) {
	quote := quote.QuoteGen()
	templates := template.Must(template.ParseFiles("templates/projects.html"))
	fmt.Println("projects")
	if err := templates.ExecuteTemplate(w, "projects.html", quote); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func statistics(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("templates/statistics.html"))
	fmt.Println("statistics")

	if err := templates.ExecuteTemplate(w, "statistics.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
