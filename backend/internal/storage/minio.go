package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"localcloud/internal/config"

	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func CreateMinioUser(accessKey, secretKey string) error {
	cfg := config.GetConfig()

	// 1. 使用 Admin Client 创建用户
	// 使用 NewWithOptions 替代弃用的 New
	adminClient, err := madmin.NewWithOptions(cfg.MinioEndpoint, &madmin.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return fmt.Errorf("admin client 初始化失败: %w", err)
	}

	// 创建 MinIO 用户
	if err := adminClient.AddUser(context.Background(), accessKey, secretKey); err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	// 2. 使用普通 Client 操作存储桶
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return fmt.Errorf("minio client 初始化失败: %w", err)
	}

	// 创建存储桶
	bucketName := fmt.Sprintf("user-%s", strings.ReplaceAll(accessKey, " ", "-"))
	if err := client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("创建存储桶失败: %w", err)
	}

	// 设置存储桶策略
	policy := fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [{
            "Effect": "Allow",
            "Principal": {"AWS": ["%s"]},
            "Action": ["s3:*"],
            "Resource": ["arn:aws:s3:::%s/*"]
        }]
    }`, accessKey, bucketName)

	if err := client.SetBucketPolicy(context.Background(), bucketName, policy); err != nil {
		return fmt.Errorf("设置策略失败: %w", err)
	}

	return nil
}

func UploadFile(userID int64, file io.Reader, size int64) (url string, err error) {
	fmt.Printf("userID: %v\n", userID)
	fmt.Printf("file: %v\n", file)
	fmt.Printf("size: %v\n", size)
	return "", nil
}
