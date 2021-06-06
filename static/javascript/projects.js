fetch("/static/json/projects_data.json")
    .then(function (response) {
        return response.json();
    })
    .then(function (data) {
        appendData(data);
    })
    .catch(function (err) {
        console.log("error: " + err);
    });

function appendData(data) {
    for (var i = 0; i < data.length; i++) {
        if (data[i].description == "") {
            continue;
        }

        var mainContainer = document.getElementById("projects");
        var project_frame = document.createElement("div");

        project_frame.setAttribute("id", data[i].name);
        project_frame.setAttribute("class", "container");

        mainContainer.prepend(project_frame);

        var subContainer = document.getElementById(data[i].name);

        var title = document.createElement("h1");
        var description = document.createElement("h2");
        var logo = document.createElement("img");
        var info = document.createElement("p");

        title.setAttribute("href", "http://localhost");
        logo.setAttribute("class", "floatl");
        logo.src =
            "https://raw.githubusercontent.com/vihdutta/autosort/master/newlogor.ico";

        title.innerHTML = data[i].name;
        description.innerHTML = data[i].description;
        info.innerHTML =
            "Size: " +
            data[i].size +
            "<br>" +
            "Stars: " +
            data[i].stargazers_count +
            "<br>" +
            "License: " +
            data[i].license.name +
            "<br><br>";

        subContainer.appendChild(title);
        subContainer.appendChild(description);
        subContainer.appendChild(logo);
        subContainer.appendChild(info);
    }
}