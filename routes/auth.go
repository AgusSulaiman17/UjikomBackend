package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func Login(router *gin.Engine) {
	router.POST("/login", controllers.Login)
}
