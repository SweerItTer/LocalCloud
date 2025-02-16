package handlers

import (
	"localcloud/internal/models"
	"localcloud/internal/storage"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	// 获取用户ID
	session := sessions.Default(c)
	user := session.Get("user").(models.User)
	userID := user.ID

	// 获取上传文件
	fileHeader, _ := c.FormFile("file")

	// 打开文件流
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "文件打开失败"})
		return
	}
	defer file.Close() // 重要：关闭文件流

	// 调用存储模块
	url, err := storage.UploadFile(
		userID,
		file,            // io.Reader
		fileHeader.Size, // int64
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "上传失败"})
		return
	}

	c.JSON(200, gin.H{"url": url})
}
