package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             int64  `gorm:"primaryKey"`
	GitHubID       int64  `gorm:"uniqueIndex"` // GitHub 用户唯一 ID，可为空
	Name           string `gorm:"not null"`    // 用户名必须有
	Email          string `gorm:"uniqueIndex"` // 允许为空，GitHub 用户可能没有
	PasswordHash   string // 允许为空，GitHub 用户没有密码
	MinioAccessKey string `gorm:"uniqueIndex"` // 允许为空，GitHub 用户可能未创建 MinIO 账号
	MinioSecretKey string // 允许为空
	AvatarURL      string `gorm:"default:'/avatars/default.png'"` // 默认头像
}
