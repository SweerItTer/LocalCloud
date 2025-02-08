document.addEventListener("DOMContentLoaded", function() {
    const avatar = document.getElementById("userAvatar");
    const loginBtn = document.getElementById("loginBtn");
    const logoutBtn = document.getElementById("logoutBtn");
    const dropdownMenu = document.getElementById("dropdownMenu");

    // **点击头像，切换下拉菜单**
    avatar.addEventListener("click", function(event) {
        event.stopPropagation(); // 防止触发 body 的点击事件
        dropdownMenu.classList.toggle("show");
    });

    // **点击页面其他部分时，关闭下拉菜单**
    document.addEventListener("click", function(event) {
        if (!dropdownMenu.contains(event.target) && !avatar.contains(event.target)) {
            dropdownMenu.classList.remove("show");
        }
    });
    
    // **检查用户是否已登录**
    function isLoggedIn() {
        return localStorage.getItem("github_user") !== null;
    }

    function getUserInfo() {
        return JSON.parse(localStorage.getItem("github_user"));
    }

    // **更新 UI**
    function updateUserUI() {
        if (isLoggedIn()) {
            const user = getUserInfo();
            avatar.src = user.avatar || "./icon/default-avatar.png";
            loginBtn.style.display = "none";
            logoutBtn.style.display = "block";
        } else {
            avatar.src = "./icon/default-avatar.png";
            loginBtn.style.display = "block";
            logoutBtn.style.display = "none";
        }
    }

    // **GitHub 登录**
    loginBtn.addEventListener("click", function() {
        console.log("🔗 跳转 GitHub 登录...");
        window.location.href = "/auth/github/login"; // **✅ 这里改为 Nginx 代理的路径**
    });

    // **处理 GitHub 回调**
    if (window.location.pathname === "/auth/github/callback") {
        const urlParams = new URLSearchParams(window.location.search);
        const code = urlParams.get("code");

        if (code) {
            console.log("🔑 GitHub 返回 code:", code);
            fetch(`/auth/github/callback?code=${code}`)
                .then(response => response.json())
                .then(user => {
                    localStorage.setItem("github_user", JSON.stringify(user));
                    console.log("✅ 登录成功，跳转主页...");
                    window.location.href = "/index.html";
                })
                .catch(error => {
                    console.error("❌ GitHub 登录失败:", error);
                    alert("GitHub 登录失败，请重试");
                    window.location.href = "/index.html";
                });
        }
    }

    // **退出登录**
    logoutBtn.addEventListener("click", function() {
        localStorage.removeItem("github_user");
        updateUserUI();
        alert("已退出 GitHub 登录");
    });

    updateUserUI();
});

