package modules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ProjectsStruct struct {
	Project     string  `json:"name"`
	Url         string  `json:"html_url"`
	Description string  `json:"description"`
	Commits     string  `json:"sha"`
	Size        int32   `json:"size"`
	Stars       int16   `json:"stargazers_count"`
	LicenseData License `json:"license"`
}

type License struct {
	License string `json:"name"`
}

func ProjectsGen() []ProjectsStruct {
	req, err := http.Get("https://api.github.com/users/vihdutta/repos?per_page=100")

	if err != nil {
		fmt.Print(err.Error())
	}

	defer req.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(req.Body)
	var projectsdata []ProjectsStruct
	json.Unmarshal(bodyBytes, &projectsdata)

	file, _ := json.MarshalIndent(projectsdata, "", " ")
	ioutil.WriteFile("static/json/projects_data.json", file, 0644)

	return projectsdata

}
