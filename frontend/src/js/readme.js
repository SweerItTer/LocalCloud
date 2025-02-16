document.addEventListener("DOMContentLoaded", function() {
    const githubReadmeUrl = "https://raw.githubusercontent.com/SweerItTer/LocalCloud/main/README.md";

    fetch(githubReadmeUrl)
        .then(response => response.ok ? response.text() : Promise.reject(`HTTP error! Status: ${response.status}`))
        .then(markdown => {
            document.getElementById("readmeText").innerHTML = marked.parse(markdown);
        })
        .catch(error => {
            document.getElementById("readmeText").innerHTML = "<p style='color: red;'>无法加载 README 文件</p>";
            console.error("Error loading README.md:", error);
        });
});
