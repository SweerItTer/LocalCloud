package main

import (
	"log"
	"net/http"

	"localcloud/internal/config"
	"localcloud/internal/handlers"
	"localcloud/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	// 优先加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, using system environment variables")
	}
	cfg := config.GetConfig()

	// 初始化 OAuth
	handlers.InitOAuth(cfg) // 将加载的配置传递给路由

	// 启动服务
	router := routes.SetupRouter(cfg)
	router.Run(cfg.ServerAddress)

	log.Printf("🚀 Server running on %s", cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, router))
}
