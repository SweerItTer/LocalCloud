document.addEventListener("DOMContentLoaded", checkLoginStatus);

async function checkLoginStatus() {
    try {
        let response = await fetch("/api/auth/check", {
            method: "GET",
            credentials: "include"  // 跨域请求时携带 Cookie
        });

        if (response.ok) {
            // 解析响应数据
            let userData = await response.json();
            console.log("用户已登录", userData);
            document.getElementById("username").innerText = userData.Name;
            // 根据需要更新页面，比如隐藏登录按钮、显示用户信息等
        } else {
            console.log("用户未登录");
        }
    } catch (error) {
        console.error("请求出错:", error);
        console.log("用户未登录");
    }
}

// GitHub 登录跳转
document.getElementById('githubLogin').addEventListener('click', () => {
    window.location.href = '/api/auth/github/login';
});

// 邮箱登录处理
document.getElementById('emailLogin').addEventListener('click', async () => {
    // 获取输入内容
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    try {
        console.error('登录请求失败:', err);
        const res = await fetch('/api/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password }), // 暂未加密(私人云盘真的需要加密吗?)
            credentials: 'include'
        });
        
        if (res.ok) {
            redirectToOriginal();
        } else {
            alert('登录失败');
        }
    } catch (err) {
        console.error('登录请求失败:', err);
    }
});

// 保存原始访问路径
function saveOriginalPath() {
    if (!window.location.pathname.includes('/login')) {
        sessionStorage.setItem('originalPath', window.location.href);
    }
}

// 登录跳转逻辑
function initLoginFlow() {
    checkLoginStatus().then(data => {
        if (!data.loggedIn) {
            saveOriginalPath();
            window.location.href = '/login.html';
        }
    });
}

// 登录后跳转
function redirectToOriginal() {
    const originalPath = sessionStorage.getItem('originalPath') || '/dashboard.html';
    sessionStorage.removeItem('originalPath');
    window.location.href = originalPath;
}