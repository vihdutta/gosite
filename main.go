package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	modules "github.com/vihdutta/gowebsite/modules"
)

type M map[string]interface{}

func main() {
	modules.ProjectsGen()
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
	dir, _ := os.Getwd()
	fmt.Println(dir + "/")

	if r.Method == "GET" {
		templates := template.Must(template.ParseFiles("templates/webapps.html"))
		if err := templates.ExecuteTemplate(w, "webapps.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	fmt.Println("webapps")

	r.ParseMultipartForm(10)
	// func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)
	file, fileHeader, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("fileHeader.Filename: %v\n", fileHeader.Filename)
	fmt.Printf("fileHeader.Size: %v\n", fileHeader.Size)
	fmt.Printf("fileHeader.Header: %v\n", fileHeader.Header)

	// tempFile, err := ioutil.TempFile("images", "upload-*.png")
	contentType := fileHeader.Header["Content-Type"][0]
	fmt.Println("Content Type:", contentType)
	var osFile *os.File
	// func TempFile(dir, pattern string) (f *os.File, err error)
	if contentType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		osFile, err = ioutil.TempFile(dir+"/"+"temp-xlsx", "*.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(dir + "/" + "temp-xlsx")
	// func ReadAll(r io.Reader) ([]byte, error)
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// func (f *File) Write(b []byte) (n int, err error)

	osFile.Write(fileBytes)
	exec.Command(dir + "/" + "zacks_requests.exe").Run()

	fmt.Println(dir + "/" + "analysis.txt")

	downloadBytes, err := ioutil.ReadFile(dir + "/" + "analysis.txt")
	if err != nil {
		fmt.Println(err)
	}

	mime := http.DetectContentType(downloadBytes)
	fileSize := len(string(downloadBytes))

	// Generate the server headers
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", "attachment; filename=analysis.txt")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", strconv.Itoa(fileSize))
	w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")

	http.ServeContent(w, r, dir+"/"+"analysis.txt", time.Now(), bytes.NewReader(downloadBytes))

	//file.Close()
	//osFile.Close()
	//defer os.Remove("analysis.txt")
	//defer os.Remove(osFile.Name())
}
