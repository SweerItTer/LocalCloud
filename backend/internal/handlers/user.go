package handlers

import (
	"localcloud/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user").(models.User)
	c.JSON(200, user)
}
