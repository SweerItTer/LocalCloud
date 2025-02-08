package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

// 用户信息结构体
type GitHubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

// 会话存储结构
type SessionStore struct {
	sync.RWMutex
	sessions map[string]*GitHubUser
}

var (
	store = &SessionStore{
		sessions: make(map[string]*GitHubUser),
	}
	cookieStore = sessions.NewCookieStore([]byte("your-secret-key"))
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// 配置路由
	http.HandleFunc("/auth/github/login", githubLoginHandler)
	http.HandleFunc("/auth/github/callback", githubCallbackHandler)
	http.HandleFunc("/user", userInfoHandler)
	http.HandleFunc("/logout", logoutHandler)

	// 启动服务器
	log.Println("✅ Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// GitHub 登录入口
func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	redirectURI := os.Getenv("REDIRECT_URI")

	authURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user",
		clientID,
		redirectURI,
	)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// GitHub 回调处理
func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	// 获取 access token
	token, err := getAccessToken(code)
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	// 获取用户信息
	user, err := getUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// 创建会话
	session, _ := cookieStore.Get(r, "user-session")
	session.Values["userID"] = strconv.Itoa(user.ID)
	session.Values["login"] = user.Login
	session.Save(r, w)

	// 存储用户信息
	store.Lock()
	store.sessions[strconv.Itoa(user.ID)] = user
	store.Unlock()

	// 重定向到前端
	http.Redirect(w, r, "/", http.StatusFound)
}

// 获取用户信息端点
func userInfoHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStore.Get(r, "user-session")
	userID, ok := session.Values["userID"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	store.RLock()
	user, exists := store.sessions[userID]
	store.RUnlock()

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// 退出登录
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStore.Get(r, "user-session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

// 辅助函数：获取 access token
func getAccessToken(code string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", nil)
	q := req.URL.Query()
	q.Add("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	q.Add("client_secret", os.Getenv("GITHUB_CLIENT_SECRET"))
	q.Add("code", code)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

// 辅助函数：获取用户信息
func getUserInfo(token string) (*GitHubUser, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Add("Authorization", "token "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
