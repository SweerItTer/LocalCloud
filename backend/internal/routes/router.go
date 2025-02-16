package routes

import (
	"localcloud/internal/config"
	"localcloud/internal/handlers"
	"localcloud/internal/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	store := cookie.NewStore([]byte(cfg.SessionSecret))        // 这里的密钥要保证稳定
	router.Use(sessions.Sessions("localcloud-session", store)) // 这里的 "session-name" 必须和获取 session 的方式一致

	// 所有API路由都挂在/api下
	api := router.Group("/api")
	{
		// 公开路由
		public := api.Group("")
		{
			public.GET("/health", handlers.HealthCheck)
			public.GET("/auth/github/login", handlers.GitHubLogin)
			public.GET("/auth/github/callback", handlers.GitHubCallback)
			public.GET("/auth/check", handlers.CheckAuth)
			public.POST("/auth/login", handlers.EmailLogin)
		}

		// 需要认证的路由
		private := api.Group("")
		private.Use(middleware.AuthRequired())
		{
			private.GET("/user", handlers.GetUser)
		}
	}

	return router
}
