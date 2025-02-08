document.addEventListener("DOMContentLoaded", function() {
    const githubReadmeUrl = "https://raw.githubusercontent.com/SweerItTer/LocalCloud/main/README.md"; // 你的 GitHub README 直链

    fetch(githubReadmeUrl)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.text();
        })
        .then(markdown => {
            // 使用 marked.js 解析 Markdown 为 HTML
            document.getElementById("readmeText").innerHTML = marked.parse(markdown);
        })
        .catch(error => {
            document.getElementById("readmeText").innerHTML = "<p>无法加载 README 文件</p>";
            console.error("Error loading README.md:", error);
        });
});
