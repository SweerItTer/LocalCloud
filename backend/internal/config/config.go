package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"localcloud/internal/models"
)

type Config struct {
	ServerAddress string // 后端服务端口(输出用,实际修改无效果)
	// minio
	MinioEndpoint  string // minio服务端口
	MinioAccessKey string // minio管理员用户名
	MinioSecretKey string // minio管理员密码
	MinioUseSSL    bool
	// GitHub 登录相关配置
	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURL  string

	SessionSecret  string
	DefaultQuotaMB int // 单次上传文件大小

	// 数据库配置
	DB         *gorm.DB // 数据库实例
	DBHost     string   // 数据库主机
	DBPort     int      // 数据库端口
	DBUser     string   // 数据库用户
	DBPassword string   // 数据库密码
	DBName     string   // 数据库名称
	DBSSLMode  string   // SSL模式
}

// 加载配置并初始化数据库
func LoadConfig() (*Config, error) {
	_ = godotenv.Load() // 忽略文件不存在的错误

	cfg := &Config{
		ServerAddress:      getEnv("SERVER_ADDR", ":8080"),
		MinioEndpoint:      getEnv("MINIO_ENDPOINT", "minio:9000"),
		MinioAccessKey:     getEnv("MINIO_ACCESS_KEY", "admin"),
		MinioSecretKey:     getEnv("MINIO_SECRET_KEY", "password123"),
		MinioUseSSL:        false,
		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GitHubRedirectURL:  getEnv("REDIRECT_URI", ""),
		SessionSecret:      getEnv("SESSION_SECRET", "supersecretkey1234567890abcdef"),
		DefaultQuotaMB:     getEnvAsInt("DEFAULT_QUOTA_MB", 1024),

		// 数据库配置
		DBHost:     getEnv("DB_HOST", "db"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "localcloud"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// 初始化数据库
	db, err := initDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("数据库初始化失败: %v", err)
	}
	cfg.DB = db
	return cfg, nil
}

// 初始化数据库连接
func initDB(cfg *Config) (*gorm.DB, error) {
	// 示例：PostgreSQL 连接
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // 打开数据库
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}
	return db, nil
}

// 读取本地.env
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetConfig 获取全局配置实例
func GetConfig() *Config {
	globalConfig, initErr := LoadConfig()
	if initErr != nil {
		panic(fmt.Sprintf("初始化配置失败: %v", initErr))
	}
	return globalConfig
}
