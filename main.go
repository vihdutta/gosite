package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	modules "github.com/vihdutta/gowebsite/modules"
)

type M map[string]interface{}

func main() {
	//modules.ProjectsGen()
	//Use this to test. REMOVE/COMMENT before pushing to Github or Heroku won't work.
	//os.Setenv("PORT", "6969")

	port := os.Getenv("PORT")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/projects", projects)
	http.HandleFunc("/statistics", statistics)
	http.HandleFunc("/webapps", webapps)
	http.Handle("/analysis.txt", http.FileServer(http.Dir("./")))

	fmt.Println("Listening on :" + port)
	http.ListenAndServe(":"+port, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("templates/index.html"))
	fmt.Println("home")

	if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func projects(w http.ResponseWriter, r *http.Request) {
	quote := modules.QuoteGen()
	/*projects := modules.Projects()*/
	templates := template.Must(template.ParseFiles("templates/projects.html"))
	fmt.Println("projects")
	if err := templates.ExecuteTemplate(w, "projects.html", M{"quote": quote, "projects": projects}); err != nil {
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

func webapps(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered webapps")
	if r.Method == "GET" {
		templates := template.Must(template.ParseFiles("templates/webapps.html"))
		if err := templates.ExecuteTemplate(w, "webapps.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	fmt.Println("Parsing Multipart Form")
	r.ParseMultipartForm(10)

	fmt.Println("Reading uploaded file's basic data")
	formFile, formFileHeader, err := r.FormFile("myFile") //reads uploaded file's basic data
	contentType := formFileHeader.Header["Content-Type"][0]
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("formFileHeader.Filename: %v\n", formFileHeader.Filename)
	fmt.Printf("formFileHeader.Size: %v\n", formFileHeader.Size)
	fmt.Printf("formFileHeader.Header: %v\n", formFileHeader.Header)
	fmt.Println("Content Type:", contentType)

	var tempFile *os.File

	fmt.Println("Creating TempFile")
	if contentType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		tempFile, err = ioutil.TempFile("temp-xlsx", "*.xlsx") //creates empty file with .xlsx extension
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Reading uploaded file's byte data")
	fileBytes, err := ioutil.ReadAll(formFile) //reads uploaded file's byte data
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Copying uploaded file's byte data to tempFile")
	tempFile.Write(fileBytes) //copies file data to the temp-file

	fmt.Println("Running Zacks Requests")
	modules.ZacksRequests(tempFile.Name())

	fmt.Println("Reading analysis.txt")
	downloadBytes, err := ioutil.ReadFile("analysis.txt")
	if err != nil {
		fmt.Println(err)
	}

	mime := http.DetectContentType(downloadBytes)
	fileSize := len(string(downloadBytes))

	// Generate the server headers
	fmt.Println("Setting headers")
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", "attachment; filename=analysis.txt")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", strconv.Itoa(fileSize))
	w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")

	fmt.Println("Downloading file")
	http.ServeContent(w, r, "analysis.txt", time.Now(), bytes.NewReader(downloadBytes))

	formFile.Close()
	tempFile.Close()
	defer os.Remove("analysis.txt")
	defer os.Remove(tempFile.Name())
}
