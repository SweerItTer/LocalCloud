package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// 初始化MinIO客户端
func initMinIO() *minio.Client {
	endpoint := "minio:9000"
	accessKey := "admin"
	secretKey := "password123"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("MinIO连接失败:", err)
	}

	// 自动创建存储桶
	ctx := context.Background()
	bucketName := "photos"
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// 如果桶已存在则忽略错误
		exists, _ := minioClient.BucketExists(ctx, bucketName)
		if !exists {
			log.Fatalln("创建存储桶失败:", err)
		}
	}

	fmt.Println("MinIO初始化完成")
	return minioClient
}

func main() {
	minioClient := initMinIO()
	r := gin.Default()

	// 文件上传接口
	r.POST("/api/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		if file == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "未选择文件"})
			return
		}

		// 打开上传文件
		fileReader, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件读取失败"})
			return
		}
		defer fileReader.Close()

		// 生成唯一文件名
		objectName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Filename)

		// 上传到MinIO
		_, err = minioClient.PutObject(
			context.Background(),
			"photos",
			objectName,
			fileReader,
			file.Size,
			minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"url": fmt.Sprintf("/images/%s", objectName),
		})
	})

	// 启动服务
	r.Run(":8080")
}