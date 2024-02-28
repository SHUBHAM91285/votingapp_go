package routes

import (
	controller "github.com/SHUBHAM91285/votingApp_go/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/user/signup", controller.SignUp())
	incomingRoutes.POST("/user/login", controller.Login())
	incomingRoutes.GET("/user/profile", controller.UserProfile())
	incomingRoutes.PATCH("/user/profile/password", controller.UpdatePassword())
}
