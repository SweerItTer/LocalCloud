package handlers

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"localcloud/internal/config"
	"localcloud/internal/models"
	"localcloud/internal/storage"
)

// GET 请求处理
var oauthConf *oauth2.Config

// 初始化 OAuth 配置
func InitOAuth(cfg *config.Config) {
	oauthConf = &oauth2.Config{
		ClientID:     cfg.GitHubClientID,
		ClientSecret: cfg.GitHubClientSecret,
		RedirectURL:  cfg.GitHubRedirectURL,
		Scopes:       []string{"user:email"}, // 请求的权限范围
		Endpoint:     githuboauth.Endpoint,
	}
}

// 检查用户登录情况
func CheckAuth(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user_id")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	c.JSON(http.StatusOK, user) // 返回用户信息
}

// GitHub 登录入口
func GitHubLogin(c *gin.Context) {
	// 生成GitHub登录链接
	url := oauthConf.AuthCodeURL("state", oauth2.AccessTypeOnline)
	c.Redirect(http.StatusTemporaryRedirect, url) // 重定向到GitHub
}

// EmailLogin 处理邮箱+密码登录
func EmailLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效请求"})
		return
	}

	db := config.GetConfig().DB

	// 查找用户
	var user models.User
	result := db.Where("email = ?", req.Email).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		// 自动注册用户
		log.Printf("自动创建用户")
		user = createNewUser(req.Email, req.Password)
		if err := db.Create(&user).Error; err != nil {
			log.Printf("用户创建失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
			return
		}

		// 创建 MinIO 用户和存储桶
		if err := storage.CreateMinioUser(user.MinioAccessKey, user.MinioSecretKey); err != nil {
			log.Printf("MinIO 用户创建失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "存储服务初始化失败"})
			return
		}
	} else if result.Error != nil {
		log.Printf("数据库查询失败: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	// 创建会话
	setUserSession(c, user)

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user": gin.H{
			"email":     user.Email,
			"avatar":    user.AvatarURL,
			"minio_key": user.MinioAccessKey, // 可选返回 MinIO 凭证
		},
	})
}

// GitHubCallback 处理 GitHub OAuth 回调
func GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少授权码"})
		return
	}

	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Token 交换失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "认证失败"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 设置超时为 10 秒
	defer cancel()                                                           // 确保在函数结束时释放资源

	client := github.NewClient(oauthConf.Client(ctx, token))

	githubUser, _, err := client.Users.Get(ctx, "")
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取用户信息"})
		return
	}

	githubID := githubUser.GetID()
	email := githubUser.GetEmail() // 可能为空
	name := githubUser.GetLogin()
	avatarURL := githubUser.GetAvatarURL()

	db := config.GetConfig().DB

	// 查找或创建 GitHub 用户
	var user models.User
	result := db.Where("git_hub_id = ?", githubID).First(&user) //通过ID查询用户
	if result.Error == gorm.ErrRecordNotFound {
		user = models.User{
			ID:             GenerateSecureUniqueID(),
			Email:          email,
			GitHubID:       githubID,
			Name:           name,
			PasswordHash:   "", // GitHub 用户无密码
			MinioAccessKey: generateMinioKey(),
			MinioSecretKey: generateMinioSecret(),
			AvatarURL:      avatarURL,
		}
		if err := db.Create(&user).Error; err != nil {
			log.Printf("GitHub 用户创建失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
			return
		}

		// 创建 MinIO 用户
		if err := storage.CreateMinioUser(user.MinioAccessKey, user.MinioSecretKey); err != nil {
			log.Printf("MinIO 用户创建失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "存储服务初始化失败"})
			return
		}
	} else if result.Error != nil {
		log.Printf("数据库查询失败: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	// 创建会话
	setUserSession(c, user)

	// 重定向到前端
	c.Redirect(http.StatusFound, "/")
}

// setUserSession 统一会话创建
func setUserSession(c *gin.Context, user models.User) {
	session := sessions.Default(c)
	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24 * 7, // 7 天有效
		HttpOnly: true,
	})
	session.Set("user_id", user.ID)
	err := session.Save() // 必须保存 session，否则不会生效
	if err != nil {
		log.Printf("保存 session 失败: %v", err)
	}

	log.Printf("Session 设定成功: user_id=%v", session.Get("user_id")) // 这行用来调试
}

func createNewUser(email, password string) models.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return models.User{
		ID:             GenerateSecureUniqueID(),
		Email:          email,
		Name:           "UserName",
		PasswordHash:   string(hashedPassword),
		MinioAccessKey: generateMinioKey(),
		MinioSecretKey: generateMinioSecret(),
		AvatarURL:      "/avatars/default-avatar.png",
	}
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

// generateMinioKey 生成一个 16 字节的随机字符串作为 MinIO Access Key
func generateMinioKey() string {
	return randomHex(16)
}

// generateMinioSecret 生成一个 32 字节的随机字符串作为 MinIO Secret Key
func generateMinioSecret() string {
	return randomHex(32)
}

// randomHex 生成指定长度的随机十六进制字符串
func randomHex(n int) string {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("无法生成随机数据: " + err.Error())
	}
	return hex.EncodeToString(bytes)
}

// 生成安全的唯一 ID（int64）
func GenerateSecureUniqueID() int64 {
	// 生成 8 字节（64 位）的随机数
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic("无法生成随机数: " + err.Error())
	}

	// 解析为 int64
	randomPart := int64(binary.LittleEndian.Uint64(b[:])) & 0x7FFFFFFFFFFFFFFF // 确保是正数

	// 获取当前毫秒级时间戳
	timestamp := time.Now().UnixMilli()

	// 组合时间戳和随机数，保证递增且唯一
	return timestamp<<16 | (randomPart & 0xFFFF) // 低 16 位用随机数填充
}
